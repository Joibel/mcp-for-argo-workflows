# Implement Single Linear Issue

Implement a specific Linear issue for mcp-for-argo-workflows.

## Usage

Provide the issue identifier (e.g., PIP-15) as an argument: `/implement-issue PIP-15`

The argument is: $ARGUMENTS

## Step 1: Fetch Issue Details

1. Use `mcp__linear-server__get_issue` with the provided issue ID
2. Parse the issue description for:
   - Tasks/requirements
   - Tool schema (if MCP tool)
   - Implementation notes
   - Dependencies
   - Acceptance criteria

## Step 2: Check Dependencies

1. Identify dependencies from the issue description (usually listed at bottom)
2. Verify each dependency is complete:
   - Check Linear status
   - Verify code exists locally
3. If dependencies not met, report and stop

## Step 3: Update Linear Status

Move issue to "In Progress":
```
mcp__linear-server__update_issue(id: "<issue-id>", state: "In Progress")
```

## Step 4: Plan Agent Collaboration

Based on issue labels and content, determine which agents need to be involved:

### Primary Implementation Agent

| Label | Primary Agent |
|-------|---------------|
| `setup` | `go-developer` or `ci-devops` |
| `mcp-tool` | `mcp-tool-implementer` |
| `testing` | `testing` |
| `docs` | `docs-examples` |
| `ci` | `ci-devops` |

### Supporting Agents (as needed)

- **`testing`** - Write or update tests for the implementation
- **`docs-examples`** - Update README, CLAUDE.md, or add examples
- **`go-developer`** - Review Go code patterns and architecture
- **`kubernetes-argo`** - Review Argo/K8s integration code

### Agent Collaboration Patterns

1. **Implementation + Testing**: Primary agent implements, then `testing` agent adds/updates tests
2. **Implementation + Docs**: Primary agent implements, then `docs-examples` updates documentation
3. **Implementation + Review**: Primary agent implements, then another agent reviews for correctness
4. **Full Pipeline**: Implement → Test → Document → Review

### Model Selection

Use the `model` parameter in the Task tool to select the appropriate model for each agent:

| Model | Use For | Examples |
|-------|---------|----------|
| `opus` | Complex architecture, deep review, critical decisions | MCP server skeleton, Argo client wrapper, cross-checking complex code |
| `sonnet` | Standard implementation tasks (default) | Individual MCP tools, CI setup, documentation, most features |

**Recommended model by task type:**

| Task Type | Primary Agent Model | Review/Cross-check Model |
|-----------|--------------------|-----------------------|
| Core architecture (PIP-10, PIP-13) | `opus` | `opus` |
| MCP tool implementation | `sonnet` | `sonnet` |
| Testing | `sonnet` | - |
| Documentation | `sonnet` | - |
| CI/DevOps | `sonnet` | `sonnet` |
| Code review / Cross-check | - | `opus` |

## Step 5: Execute with Multiple Agents

### Phase 1: Primary Implementation

Delegate to the primary agent with:
- Full issue description
- Any relevant context from existing code
- Clear instruction to implement according to spec
- **Model**: `opus` for core architecture, `sonnet` for standard tasks

### Phase 2: Supporting Agents (run as appropriate)

After primary implementation, invoke supporting agents:

**If code was written, consider:**
- `testing` agent (model: `sonnet`): "Review the implementation of [issue] and add appropriate unit tests"
- `go-developer` agent (model: `sonnet`): "Review the implementation of [issue] for Go best practices and potential issues"

**If a new feature was added, consider:**
- `docs-examples` agent (model: `sonnet`): "Update documentation to reflect the new [feature] implementation"

**If MCP tool was implemented, consider:**
- `testing` agent (model: `sonnet`): "Add unit tests for the new [tool_name] MCP tool"
- `docs-examples` agent (model: `sonnet`): "Add usage example for the new [tool_name] tool to the README"

### Phase 3: Cross-Check (optional but recommended)

For complex implementations, have a different agent verify using `opus` for deeper analysis:
- `go-developer` (model: `opus`): "Review this implementation for correctness, error handling, and edge cases"
- `kubernetes-argo` (model: `opus`): "Verify the Argo client usage is correct and follows best practices"

### Phase 4: Create Follow-up Tasks (as needed)

During implementation, agents may identify improvements or issues that are out of scope for the current task. Instead of expanding scope, create new Linear issues for later:

**When to create follow-up tasks:**
- Technical debt that should be addressed later
- Refactoring opportunities discovered during implementation
- Additional test coverage needed
- Documentation improvements
- Performance optimizations
- Edge cases that need handling
- Related features that would be nice to have

**How to create a follow-up task:**
```
mcp__linear-server__create_issue(
  team: "Pipekit",
  project: "mcp-for-argo-workflows",
  title: "Brief description of the task",
  description: "## Context\n\nDiscovered while implementing [PIP-X].\n\n## Problem/Opportunity\n\n[Description]\n\n## Suggested Approach\n\n[How to fix/improve]\n\n## Dependencies\n\n- PIP-X (if applicable)",
  labels: ["technical-debt"] or ["enhancement"] or ["testing"] as appropriate
)
```

**Guidelines for follow-up tasks:**
- Keep the current task focused - don't expand scope
- Be specific about what needs to be done
- Reference the original issue for context
- Use appropriate labels: `technical-debt`, `enhancement`, `testing`, `docs`, `bug`
- Don't create follow-ups for trivial issues - fix them now if quick

## Step 6: Verify Implementation

1. **Run linter** - `make lint` (fix any issues)
2. **Run tests** - `make test` (ensure passing)
3. **Manual check** - Verify against acceptance criteria from issue

## Step 7: Commit Changes

Create a commit with message format:
```
[PIP-X] Brief description of implementation

- Detail 1
- Detail 2

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Step 8: Update Linear

1. Move issue to "Done":
   ```
   mcp__linear-server__update_issue(id: "<issue-id>", state: "Done")
   ```

2. Add completion comment listing what was done:
   ```
   mcp__linear-server__create_comment(issueId: "<issue-id>", body: "Implementation complete. Commit: <hash>\n\nAgents involved:\n- [primary agent]: [what they did]\n- [supporting agent]: [what they did]\n\nFollow-up tasks created:\n- PIP-XX: [title] (if any)")
   ```

## Error Handling

If implementation fails:
1. Keep issue in "In Progress"
2. Add comment describing the blocker
3. Report to user with details
4. Suggest resolution steps

## Examples

### Example 1: MCP Tool Implementation (PIP-15)

1. **Primary**: `mcp-tool-implementer` (model: `sonnet`) - Implements submit_workflow tool
2. **Supporting**: `testing` (model: `sonnet`) - Adds unit tests for the tool handler
3. **Supporting**: `docs-examples` (model: `sonnet`) - Updates README with tool description

### Example 2: Core Architecture (PIP-10)

1. **Primary**: `go-developer` (model: `opus`) - Implements MCP server skeleton
2. **Supporting**: `testing` (model: `sonnet`) - Adds basic server tests
3. **Cross-check**: `mcp-tool-implementer` (model: `opus`) - Verifies tool registration pattern is correct

### Example 3: CI Setup (PIP-8)

1. **Primary**: `ci-devops` (model: `sonnet`) - Creates GitHub Actions workflow
2. **Cross-check**: `go-developer` (model: `sonnet`) - Verifies Go-specific CI configuration

---

Begin by fetching the issue details for: $ARGUMENTS
