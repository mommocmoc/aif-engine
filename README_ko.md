# AIF (Agent Interfaces Framework) 🤖💻

> **무거운 MCP 서버 구축은 멈추세요. AI 에이전트가 진짜 원하는 건 빠르고, 조합 가능한 단일 바이너리 CLI입니다.**

[Read in English](README.md)

## 문제점 (The Problem)
현재 AI 에이전트 생태계는 MCP(Model Context Protocol)에 지나치게 집착하고 있습니다. 물론 MCP가 훌륭한 도구이긴 하지만, 이를 위해선 무거운 백그라운드 데몬(서버)을 돌리고, Node.js나 Python 의존성을 관리해야 하며, 복잡한 네트워크 프로토콜을 다뤄야 합니다.

AI 에이전트가 로컬 컴퓨터에서 작업을 수행할 때, 그들은 느린 웹 서버와 통신하고 싶어 하지 않습니다. 그들은 **터미널**을 원하고, 파이프라인(`|`)으로 결과를 연결하고 싶어 합니다. 즉, **유닉스 철학(UNIX philosophy)**을 원합니다.

## 해결책: AIF
**AIF (Agent Interfaces Framework)**는 단순한 JSON/YAML 형태의 API 명세서(Specification)만 주어지면, AI 친화적인 CLI 도구를 자동으로 생성해 주는 오픈소스 엔진입니다.

모든 SaaS 도구마다 일일이 커스텀 MCP 서버를 작성하는 대신, AIF는 Claude Code, OpenClaw, GitHub Copilot 같은 에이전트들이 즉시 사용할 수 있도록 100% 순수 JSON만을 출력하는 가벼운 Go 언어 기반 CLI를 생성하는 것을 지향합니다.

### AI 에이전트의 새로운 워크플로우
1. AI 에이전트가 특정 서비스의 API 문서를 읽습니다.
2. 에이전트가 엔드포인트와 플래그를 정의한 간단한 `spec.json` 파일을 작성합니다.
3. 에이전트가 `aif build spec.json` 명령어를 실행합니다.
4. **자동 토큰 인증(`auth login`) 기능과 JSON 포맷 출력이 완벽하게 구현된 Go 바이너리 CLI가 단 몇 초 만에 생성됩니다.**

## 개념 증명 (Proof of Concept)
이 저장소에는 AIF 엔진의 PoC 코드와 함께, [Upload-Post API](https://upload-post.com)를 제어하기 위한 예시 명세서인 `upost-spec.json`이 포함되어 있습니다.

### 사용해 보기:
```bash
# AIF 엔진 빌드
go build -o aif main.go

# 명세서(spec)로부터 타겟 CLI(예: upost) 자동 생성
./aif build upost-spec.json

# 생성된 CLI 테스트
./upost help
./upost auth login --token "당신의-API-키"
./upost text --title "Hello World" --platform x --platform linkedin
```

## 함께 프로젝트를 키워나갈 분들을 찾습니다
*이 프로젝트는 AI가 내 컴퓨터에서 더 완벽하게 일하길 바랐던 한 비개발자의 아이디어에서 출발했으며, AI와 함께 코드를 작성했습니다.*

AI 시대에 걸맞은 빠르고(Fast), 상태를 유지하지 않으며(Stateless), 결합 가능한(Composable) 도구를 만드는 데 공감하는 Go 및 Rust 개발자분들의 많은 관심과 참여(Issue, PR)를 기다립니다!