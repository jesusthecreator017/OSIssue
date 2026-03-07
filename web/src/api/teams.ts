import { apiFetch } from "./client";
import type {
	Team,
	TeamMember,
	CreateTeamInput,
	AddTeamMemberInput,
} from "../schemas/teams";

export const teamsApi = {
	list: () => apiFetch<{ teams: Team[] }>("/v1/teams").then((r) => r.teams),

	getById: (id: string) =>
		apiFetch<{ team: Team }>(`/v1/teams/${id}`).then((r) => r.team),

	create: (data: CreateTeamInput) =>
		apiFetch<{ team: Team }>("/v1/teams", {
			method: "POST",
			body: JSON.stringify(data),
		}).then((r) => r.team),

	delete: (id: string) =>
		apiFetch<{ message: string }>(`/v1/teams/${id}`, {
			method: "DELETE",
		}),

	getMembers: (teamId: string) =>
		apiFetch<{ members: TeamMember[] }>(`/v1/teams/${teamId}/members`).then(
			(r) => r.members,
		),

	addMember: (teamId: string, data: AddTeamMemberInput) =>
		apiFetch<{ team_member: TeamMember }>(`/v1/teams/${teamId}/members`, {
			method: "POST",
			body: JSON.stringify(data),
		}).then((r) => r.team_member),

	removeMember: (teamId: string, userId: string) =>
		apiFetch<{ message: string }>(`/v1/teams/${teamId}/members/${userId}`, {
			method: "DELETE",
		}),
};
