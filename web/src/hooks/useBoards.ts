"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { boardsApi } from "@/api/boards";

export function usePersonalBoard() {
	return useQuery({
		queryKey: ["personal-board"],
		queryFn: boardsApi.getPersonalBoard,
	});
}

export function useTeamBoard(teamId: string) {
	return useQuery({
		queryKey: ["team-board", teamId],
		queryFn: () => boardsApi.getTeamBoard(teamId),
		enabled: !!teamId,
	});
}

export function useBoardIssues(boardId: string | undefined) {
	return useQuery({
		queryKey: ["board-issues", boardId],
		queryFn: () => boardsApi.getBoardIssues(boardId as string),
		enabled: !!boardId,
	});
}

export function useCreateColumn() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({
			boardId,
			name,
			position,
		}: {
			boardId: string;
			name: string;
			position: number;
		}) => boardsApi.createColumn(boardId, name, position),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["personal-board"] });
			qc.invalidateQueries({ queryKey: ["team-board"] });
		},
	});
}

export function useUpdateColumn() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({
			boardId,
			colId,
			name,
		}: {
			boardId: string;
			colId: string;
			name: string;
		}) => boardsApi.updateColumn(boardId, colId, name),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["personal-board"] });
			qc.invalidateQueries({ queryKey: ["team-board"] });
		},
	});
}

export function useReorderColumn() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({
			boardId,
			colId,
			position,
		}: {
			boardId: string;
			colId: string;
			position: number;
		}) => boardsApi.reorderColumn(boardId, colId, position),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["personal-board"] });
			qc.invalidateQueries({ queryKey: ["team-board"] });
		},
	});
}

export function useDeleteColumn() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: ({ boardId, colId }: { boardId: string; colId: string }) =>
			boardsApi.deleteColumn(boardId, colId),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["personal-board"] });
			qc.invalidateQueries({ queryKey: ["team-board"] });
		},
	});
}
