import { z } from "zod";

export const TeamRole = z.enum(["owner", "admin", "member"]);
export type TeamRole = z.infer<typeof TeamRole>;

export const TeamSchema = z.object({
	id: z.uuid(),
	name: z.string(),
	description: z.string().optional(),
	created_by: z.uuid(),
	max_members: z.number().min(5, "Minimum Max members is 5").optional(),
	avatar_url: z.string().optional(),
	created_at: z.string(),
	updated_at: z.string(),
});
export type Team = z.infer<typeof TeamSchema>;

export const TeamMemberSchema = z.object({
	user_id: z.uuid(),
	team_id: z.uuid(),
	user_name: z.string().optional(),
	email: z.string().optional(),
	role: TeamRole,
	joined_at: z.string(),
});
export type TeamMember = z.infer<typeof TeamMemberSchema>;

export const ListUserTeamsRowSchema = z.object({
	id: z.uuid(),
	name: z.string(),
	description: z.string(),
	role: TeamRole,
	joined_at: z.string(),
	created_at: z.string(),
	updated_at: z.string(),
});
export type ListUserTeamsRow = z.infer<typeof ListUserTeamsRowSchema>;

export const CreateTeamSchema = z.object({
	name: z.string().min(1, "Name is required").max(127, "Name too long"),
	description: z.string().max(500, "Description too long").default(""),
	max_members: z.number().min(1).optional(),
	avatar_url: z.string().optional(),
});
export type CreateTeamInput = z.infer<typeof CreateTeamSchema>;

export const AddTeamMemberSchema = z.object({
	user_id: z.uuid(),
	role: TeamRole,
});
export type AddTeamMemberInput = z.infer<typeof AddTeamMemberSchema>;
