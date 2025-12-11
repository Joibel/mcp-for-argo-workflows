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

## Step 3: Select Specialist

Based on issue labels and content:

| Label | Agent |
|-------|-------|
| `setup` | `go-developer` or `ci-devops` |
| `mcp-tool` | `mcp-tool-implementer` |
| `testing` | `testing` |
| `docs` | `docs-examples` |
| `ci` | `ci-devops` |

## Step 4: Update Linear Status

Move issue to "In Progress":
```
mcp__linear-server__update_issue(id: "<issue-id>", state: "In Progress")
```

## Step 5: Implement

Delegate to the appropriate specialist agent with:
- Full issue description
- Any relevant context from existing code
- Clear instruction to implement according to spec

## Step 6: Verify Implementation

1. **Code review** - Check implementation matches requirements
2. **Run linter** - `make lint` (fix any issues)
3. **Run tests** - `make test` (ensure passing)
4. **Manual check** - Verify against acceptance criteria

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

2. Add completion comment:
   ```
   mcp__linear-server__create_comment(issueId: "<issue-id>", body: "Implementation complete. Commit: <hash>")
   ```

## Error Handling

If implementation fails:
1. Keep issue in "In Progress"
2. Add comment describing the blocker
3. Report to user with details
4. Suggest resolution steps

---

Begin by fetching the issue details for: $ARGUMENTS
