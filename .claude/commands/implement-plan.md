# Implement Plan from Linear

You are orchestrating the implementation of the mcp-for-argo-workflows project based on the Linear project plan.

## Step 1: Fetch Current State from Linear

First, retrieve the current project status:

1. Use `mcp__linear-server__list_issues` with `project: "mcp-for-argo-workflows"` to get all issues
2. Identify issues by status: Backlog, Todo, In Progress, Done
3. Check issue dependencies (noted in each issue's description)

## Step 2: Identify Ready Tasks

Analyze the issues to find tasks that are ready to implement:

1. **Setup tasks** (PIP-5 through PIP-14) - Foundation work, check dependencies
2. **MCP tools** (PIP-15 through PIP-49) - Require setup tasks complete
3. **Testing** (PIP-50, PIP-51) - Require tools implemented
4. **Documentation** (PIP-9, PIP-52, PIP-53) - Can be done incrementally

A task is "ready" when:
- Status is Backlog or Todo
- All dependencies (listed in issue description) are completed

## Step 3: Create Implementation Plan

For each ready task, determine:

1. **Which specialist agent** should handle it:
   - `go-developer` - Core Go code, server skeleton, client wrapper
   - `mcp-tool-implementer` - Individual MCP tool implementations
   - `kubernetes-argo` - Argo client logic, K8s integration
   - `testing` - Unit tests, mocks, integration tests
   - `ci-devops` - GitHub Actions, Makefile, linting config
   - `docs-examples` - README, examples, documentation

2. **Implementation order** based on dependencies

3. **Verification criteria** from the issue description

## Step 4: Execute Implementation

For each task in order:

1. **Update Linear** - Move issue to "In Progress" using `mcp__linear-server__update_issue`
2. **Delegate to specialist** - Use the Task tool with appropriate subagent
3. **Verify output** - Check that implementation meets issue requirements
4. **Run tests** - Execute `make test` and `make lint` if applicable
5. **Update Linear** - Move issue to "Done" and add completion comment

## Step 5: Report Progress

After each implementation cycle:

1. Summarize what was completed
2. List any blockers or issues encountered
3. Identify next tasks ready for implementation
4. Update Linear with progress comments using `mcp__linear-server__create_comment`

## Implementation Guidelines

- **One task at a time** - Complete and verify before moving to next
- **Commit frequently** - After each logical unit of work
- **Update Linear** - Keep issue status current
- **Ask for clarification** - If requirements are ambiguous, ask before implementing
- **Run verification** - Always run `make lint` and `make test` before marking complete

## Starting Point

If this is the first implementation session, start with:
1. PIP-5: Initialize Go module and directory structure
2. PIP-6: Create Makefile
3. PIP-7: Configure golangci-lint

These have no dependencies and enable all subsequent work.

---

Begin by fetching the current Linear project state and identifying ready tasks.
