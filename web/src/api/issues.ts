import { apiFetch } from "./client";
import type {
	Issue,
	CreateIssueInput,
	UpdateIssueInput,
	MoveIssueInput,
	Label,
} from "../schemas/issue";

export const issuesApi = {
	list: () => apiFetch<{ issues: Issue[] }>("/v1/issues").then((r) => r.issues),

	getById: (id: string) =>
		apiFetch<{ issue: Issue }>(`/v1/issues/${id}`).then((r) => r.issue),

	create: (data: CreateIssueInput) =>
		apiFetch<{ issue: Issue }>("/v1/issues", {
			method: "POST",
			body: JSON.stringify(data),
		}).then((r) => r.issue),

	update: (id: string, data: UpdateIssueInput) =>
		apiFetch<{ issue: Issue }>(`/v1/issues/${id}`, {
			method: "PATCH",
			body: JSON.stringify(data),
		}).then((r) => r.issue),

	move: (id: string, data: MoveIssueInput) =>
		apiFetch<{ message: string }>(`/v1/issues/${id}/move`, {
			method: "PATCH",
			body: JSON.stringify(data),
		}),

	delete: (id: string) =>
		apiFetch<{ message: string }>(`/v1/issues/${id}`, {
			method: "DELETE",
		}),

	listLabels: (id: string) =>
		apiFetch<{ labels: Label[] }>(`/v1/issues/${id}/labels`).then(
			(r) => r.labels,
		),

	addLabel: (issueId: string, labelId: string) =>
		apiFetch<{ message: string }>(`/v1/issues/${issueId}/labels`, {
			method: "POST",
			body: JSON.stringify({ label_id: labelId }),
		}),

	removeLabel: (issueId: string, labelId: string) =>
		apiFetch<{ message: string }>(`/v1/issues/${issueId}/labels/${labelId}`, {
			method: "DELETE",
		}),
};
