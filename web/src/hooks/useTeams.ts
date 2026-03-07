"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { teamsApi } from "@/api/teams";
import type { CreateTeamInput, AddTeamMemberInput } from "@/schemas/teams";

export function useTeams() {
	return useQuery({
		queryKey: ["teams"],
		queryFn: teamsApi.list,
	});
}

export function useTeam(id: string) {
	return useQuery({
		queryKey: ["teams", id],
		queryFn: () => teamsApi.getById(id),
		enabled: !!id,
	});
}

export function useTeamMembers(teamId: string) {
	return useQuery({
		queryKey: ["teams", teamId, "members"],
		queryFn: () => teamsApi.getMembers(teamId),
		enabled: !!teamId,
	});
}

export function useCreateTeam() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (data: CreateTeamInput) => teamsApi.create(data),
		onSuccess: () => qc.invalidateQueries({ queryKey: ["teams"] }),
	});
}

export function useDeleteTeam() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (id: string) => teamsApi.delete(id),
		onSuccess: () => qc.invalidateQueries({ queryKey: ["teams"] }),
	});
}

export function useAddTeamMember(teamId: string) {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (data: AddTeamMemberInput) => teamsApi.addMember(teamId, data),
		onSuccess: () =>
			qc.invalidateQueries({ queryKey: ["teams", teamId, "members"] }),
	});
}

export function useRemoveTeamMember(teamId: string) {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (userId: string) => teamsApi.removeMember(teamId, userId),
		onSuccess: () =>
			qc.invalidateQueries({ queryKey: ["teams", teamId, "members"] }),
	});
}
