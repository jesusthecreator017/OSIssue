"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { labelsApi } from "@/api/labels";
import { issuesApi } from "@/api/issues";

export function useLabels() {
	return useQuery({
		queryKey: ["labels"],
		queryFn: labelsApi.list,
	});
}

export function useCreateLabel() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({ name, color }: { name: string; color: string }) =>
			labelsApi.create(name, color),
		onSuccess: () => qc.invalidateQueries({ queryKey: ["labels"] }),
	});
}

export function useIssueLabels(issueId: string) {
	return useQuery({
		queryKey: ["issue-labels", issueId],
		queryFn: () => issuesApi.listLabels(issueId),
		enabled: !!issueId,
	});
}

export function useAddLabelToIssue() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({ issueId, labelId }: { issueId: string; labelId: string }) =>
			issuesApi.addLabel(issueId, labelId),
		onSuccess: () => qc.invalidateQueries({ queryKey: ["issue-labels"] }),
	});
}

export function useRemoveLabelFromIssue() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({ issueId, labelId }: { issueId: string; labelId: string }) =>
			issuesApi.removeLabel(issueId, labelId),
		onSuccess: () => qc.invalidateQueries({ queryKey: ["issue-labels"] }),
	});
}
