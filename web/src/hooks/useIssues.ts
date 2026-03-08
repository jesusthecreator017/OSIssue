"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { issuesApi } from "@/api/issues";
import type {
	CreateIssueInput,
	UpdateIssueInput,
	MoveIssueInput,
} from "@/schemas/issue";

export function useIssues() {
	return useQuery({
		queryKey: ["issues"],
		queryFn: issuesApi.list,
	});
}

export function useIssue(id: string) {
	return useQuery({
		queryKey: ["issues", id],
		queryFn: () => issuesApi.getById(id),
		enabled: !!id,
	});
}

export function useCreateIssue() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (data: CreateIssueInput) => issuesApi.create(data),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["issues"] });
			qc.invalidateQueries({ queryKey: ["board-issues"] });
		},
	});
}

export function useUpdateIssue() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: string; data: UpdateIssueInput }) =>
			issuesApi.update(id, data),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["issues"] });
			qc.invalidateQueries({ queryKey: ["board-issues"] });
		},
	});
}

export function useMoveIssue() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: string; data: MoveIssueInput }) =>
			issuesApi.move(id, data),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["board-issues"] });
		},
	});
}

export function useDeleteIssue() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (id: string) => issuesApi.delete(id),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["issues"] });
			qc.invalidateQueries({ queryKey: ["board-issues"] });
		},
	});
}
