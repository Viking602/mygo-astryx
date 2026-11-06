#!/usr/bin/env python3
"""Verify the built native Gallery through MyGo's opt-in local automation API.

Run after: go build -tags e2e -o bin/gallery-e2e ./cmd/gallery
Only the child process started here is terminated during cleanup.
"""
import base64
from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path
import socket
import subprocess
import tempfile
import time

root = Path(__file__).resolve().parents[1]
artifacts = root / "artifacts"
artifacts.mkdir(exist_ok=True)
index = json.loads(subprocess.check_output([str(root / "bin/gallery-e2e"), "-list"], cwd=root))
demo_by_component = {name: demo for demo in index["Demos"] for name in demo["Components"]}
receipt = artifacts / "native-smoke.json"
receipt.write_text(json.dumps({"passed": False, "status": "running"}, indent=2) + "\n")

with tempfile.TemporaryDirectory(prefix="mygo-astryx-e2e-", dir="/tmp") as tmp:
    socket_path = str(Path(tmp) / "gallery.sock")
    env = dict(os.environ, MYGO_AUTOMATION=socket_path)
    with (artifacts / "native-gallery.log").open("w") as log:
        process = subprocess.Popen([str(root / "bin/gallery-e2e"), "-theme", "light", "-focus"], cwd=root, env=env, stdout=log, stderr=log)
        connection = None
        try:
            deadline = time.monotonic() + 15
            while time.monotonic() < deadline:
                if process.poll() is not None:
                    raise RuntimeError("Gallery 在创建窗口前退出；查看 artifacts/native-gallery.log")
                candidate = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
                try:
                    candidate.connect(socket_path)
                    connection = candidate
                    break
                except (FileNotFoundError, ConnectionRefusedError):
                    candidate.close()
                    time.sleep(0.1)
            if connection is None:
                raise RuntimeError("Gallery 自动化连接未就绪；查看 artifacts/native-gallery.log")
            connection.settimeout(15)
            reader = connection.makefile("rb")
            request_id = 0

            def call(method, params=None):
                global request_id
                request_id += 1
                payload = {"id": request_id, "method": method}
                if params is not None:
                    payload["params"] = params
                connection.sendall((json.dumps(payload) + "\n").encode())
                answer = json.loads(reader.readline())
                if "error" in answer:
                    raise RuntimeError(f"{method}: {answer['error']}")
                return answer["result"]

            windows = call("windows")["windows"]
            if not windows:
                raise RuntimeError("Gallery 没有创建窗口")
            window = windows[0]
            assert window["kind"] == "content", window
            assert window["visible"] and "MyGo / Astryx" in window["title"], window
            wid = window["id"]

            def tree():
                return call("tree", {"window": wid})

            def find(label, role=None, action=None):
                for _ in range(30):
                    # Floating controls follow their background counterparts.
                    # Match an actionable node, rather than its text child.
                    for node in reversed(tree()["nodes"]):
                        if node.get("label") == label and (role is None or node.get("role") == role) and (action is None or action in node.get("actions", [])):
                            return node
                    time.sleep(0.05)
                raise AssertionError(f"找不到原生元素：{label}")

            def act(label, action="press", role=None, **extra):
                node = find(label, role, action)
                return call("act", {"window": wid, "id": node["id"], "action": action, **extra})

            def capture(name):
                shot = call("capture", {"window": wid})
                (artifacts / name).write_bytes(base64.b64decode(shot["png"]))
                return {"width": shot["width"], "height": shot["height"]}

            def select_component(name):
                act("搜索组件", "setValue", role="textField", text=name)
                act("查看组件", role="popUpButton")
                act(name)
                for _ in range(30):
                    picker = find("查看组件", "popUpButton")
                    if picker.get("value") == name and "expanded" not in picker.get("states", []):
                        break
                    time.sleep(0.05)
                else:
                    raise AssertionError(f"组件选择未提交：{name}，{picker}")
                find(demo_by_component[name]["Title"])

            call("idle", {"window": wid, "timeout": 5000})
            size = capture("gallery-native-light.png")
            act("浏览组件", role="button")
            act("主要动作", role="button")
            find("主要动作和添加已触发 1 次", "text")
            act("深色", role="radio")
            act("浮层与对话", role="button")
            act("打开对话框", role="button")
            find("编辑项目", "dialog")
            act("名称", "setValue", role="textField", text="原生验证")
            act("方案", role="popUpButton")
            find("团队方案")
            call("input", {"window": wid, "events": [{"type": "key", "key": "Escape"}]})
            find("编辑项目", "dialog")
            capture("gallery-native-dialog-dark.png")
            call("input", {"window": wid, "events": [{"type": "key", "key": "Escape"}]})
            call("idle", {"window": wid, "timeout": 5000})
            assert not any(n.get("role") == "dialog" for n in tree()["nodes"])

            # Traverse the actual native controls for every public entry. Parts
            # intentionally share a live demo, while each selector path is used.
            components, demos = [], set()
            native_dir = artifacts / "native-components"
            native_dir.mkdir(exist_ok=True)
            for component in index["Components"]:
                name = component["Name"]
                select_component(name)
                nodes = tree()["nodes"]
                assert not any(n.get("role") == "webView" for n in nodes), name
                demo = demo_by_component[name]
                components.append({"component": name, "demo": demo["ID"], "native_nodes": len(nodes)})
                if demo["ID"] not in demos:
                    capture("native-components/" + demo["ID"] + "-dark.png")
                    demos.add(demo["ID"])
                if len(components) % 25 == 0:
                    print(f"已验证 {len(components)} 个原生组件入口", flush=True)

            select_component("MultiSelector")
            act("语言多选", role="popUpButton")
            act("筛选语言多选", "setValue", role="textField", text="Rust")
            act("Rust", role="checkBox")
            act("完成", role="button")
            find("已选择：go, rust")

            select_component("DropdownMenu")
            act("项目操作")
            act("自动保存")
            find("动作： · 自动保存：true · stable")

            select_component("CommandPalette")
            act("打开命令面板", role="button")
            act("搜索命令", "setValue", role="textField", text="设置")
            call("input", {"window": wid, "events": [{"type": "key", "key": "Enter"}]})
            find("已执行：打开设置")

            select_component("AvatarStatusDot")
            find("Ada，在线", "image")
            find("Grace，离线", "image")
            capture("gallery-native-avatar-status-dark.png")
            select_component("AvatarGroupOverflow")
            act("另外 2 人", role="button")
            find("协作成员名单", "dialog")
            call("input", {"window": wid, "events": [{"type": "key", "key": "Escape"}]})

            select_component("MobileNavToggle")
            act("移动工作空间")
            act("移动项目")
            find("移动当前位置：projects")
            act("打开站点导航", role="button")
            find("移动站点导航", "dialog")
            act("移动设置", role="button")
            find("移动当前位置：settings")

            select_component("AlertDialog")
            act("显示长确认说明", role="switch")
            act("危险操作确认", role="button")
            find("删除项目", "alertDialog")
            cancel = find("取消", "button")
            box = cancel["bounds"]
            assert box["height"] > 0 and box["y"] >= 0 and box["y"] + box["height"] <= window["contentBounds"]["height"], box
            capture("gallery-native-long-confirm-dark.png")
            act("取消", role="button")

            select_component("TopNavMegaMenuItem")
            act("显示长菜单内容", role="switch")
            act("探索产品", "scrollIntoView")
            act("探索产品")
            act("扩展示例 24", "scrollIntoView", role="button")
            capture("gallery-native-mega-scroll-dark.png")
            act("扩展示例 24", role="button")
            find("导航动作：扩展示例 24")
            assert not any(n.get("label") == "扩展示例 24" for n in tree()["nodes"])

            select_component("DateRangeInput")
            act("排除周末", role="checkBox")
            find("日期范围包含不可用日期")
            act("清除排期", role="button")
            act("工作日一天", role="button")
            assert not any(n.get("label") == "日期范围包含不可用日期" for n in tree()["nodes"])

            select_component("BottomSheet")
            act("分步底部面板", role="button")
            act("流程项目名称", "setValue", role="textField", text="原生分步验证")
            act("下一步：核对", role="button")
            find("核对资料", "dialog")
            find("项目名称：原生分步验证")
            capture("gallery-native-sheet-review-dark.png")
            call("input", {"window": wid, "events": [{"type": "key", "key": "Escape"}]})

            select_component("Lightbox")
            act("晨光", role="button")
            find("图片预览", role="dialog")
            act("放大", role="button")
            find("1 / 3 · 晨光 · 125%")
            act("下一张", role="button")
            find("2 / 3 · 夜航 · 100%")
            capture("gallery-native-lightbox-dark.png")
            call("input", {"window": wid, "events": [{"type": "key", "key": "Escape"}]})

            select_component("Chat")
            act("消息", "setValue", role="textField", text="原生聊天验证")
            act("消息", "focus", role="textField")
            call("input", {"window": wid, "events": [{"type": "key", "key": "Enter"}]})
            find("原生聊天验证")
            assert find("消息", "textField").get("value", "") == ""
            find("@项目", "group")
            capture("gallery-native-chat-dark.png")

            act("搜索组件", "setValue", role="textField", text="CheckboxInput")
            find("开关与选项")
            act("Astryx 组件对照")
            find("Astryx · CheckboxInput")
            find("查看 CheckboxInput 源码", "link")
            capture("gallery-native-astryx-reference-dark.png")

            report = {"passed": True, "completed_utc": datetime.now(timezone.utc).isoformat(), "binary_sha256": hashlib.sha256((root / "bin/gallery-e2e").read_bytes()).hexdigest(), "astryx_source": index["AstryxSource"], "window": window, "capture": size, "component_entries": len(components), "live_demos": len(demos), "components": components, "checks": ["native-window", "button-action", "theme-switch", "dialog", "text-edit", "nested-select", "escape-order", "all-component-selectors", "native-only-trees", "searchable-multi-selector", "checked-menu", "command-keyboard", "avatar-status-semantics", "avatar-overflow-action", "mobile-nested-navigation", "mobile-drawer-selection", "long-alert-fixed-actions", "mega-menu-scroll-select-dismiss", "date-rule-presets", "controlled-sheet-flow", "lightbox-zoom-navigation", "chat-enter-send", "chat-context-preserved", "astryx-alias-reference", "capture"]}
            receipt.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n")
            print(json.dumps({k: v for k, v in report.items() if k != "components"}, ensure_ascii=False))
        except Exception as error:
            receipt.write_text(json.dumps({"passed": False, "error": str(error),"completed_utc":datetime.now(timezone.utc).isoformat(),"binary_sha256":hashlib.sha256((root/"bin/gallery-e2e").read_bytes()).hexdigest(),"astryx_source":index["AstryxSource"]}, ensure_ascii=False, indent=2) + "\n")
            try:
                (artifacts / "native-failure-tree.json").write_text(json.dumps(tree(),ensure_ascii=False,indent=2)+"\n")
                capture("native-failure.png")
            except Exception:
                pass
            raise
        finally:
            if connection is not None:
                connection.close()
            if process.poll() is None:
                process.terminate()
                try:
                    process.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    process.kill()
                    process.wait()
