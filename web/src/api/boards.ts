import { apiFetch } from "./client";
import type { Board, BoardColumn, Issue } from "../schemas/issue";

export const boardsApi = {
	getPersonalBoard: () =>
		apiFetch<{ board: Board; columns: BoardColumn[] }>("/v1/me/board"),

	getTeamBoard: (teamId: string) =>
		apiFetch<{ board: Board; columns: BoardColumn[] }>(
			`/v1/teams/${teamId}/board`,
		),

	getBoardIssues: (boardId: string) =>
		apiFetch<{ issues: Issue[] }>(`/v1/boards/${boardId}/issues`).then(
			(r) => r.issues,
		),

	createColumn: (boardId: string, name: string, position: number) =>
		apiFetch<{ column: BoardColumn }>(`/v1/boards/${boardId}/columns`, {
			method: "POST",
			body: JSON.stringify({ name, position }),
		}).then((r) => r.column),

	updateColumn: (boardId: string, colId: string, name: string) =>
		apiFetch<{ column: BoardColumn }>(
			`/v1/boards/${boardId}/columns/${colId}`,
			{
				method: "PATCH",
				body: JSON.stringify({ name }),
			},
		).then((r) => r.column),

	reorderColumn: (boardId: string, colId: string, position: number) =>
		apiFetch<{ column: BoardColumn }>(
			`/v1/boards/${boardId}/columns/${colId}/reorder`,
			{
				method: "PATCH",
				body: JSON.stringify({ position }),
			},
		).then((r) => r.column),

	deleteColumn: (boardId: string, colId: string) =>
		apiFetch<{ message: string }>(`/v1/boards/${boardId}/columns/${colId}`, {
			method: "DELETE",
		}),
};
