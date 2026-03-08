import { z } from "zod";

export const PriorityType = z.enum(["Low", "Medium", "High", "Critical"]);
export type PriorityType = z.infer<typeof PriorityType>;

export const LabelSchema = z.object({
	id: z.uuid(),
	name: z.string(),
	color: z.string(),
	created_at: z.string(),
});
export type Label = z.infer<typeof LabelSchema>;

export const BoardColumnSchema = z.object({
	id: z.uuid(),
	board_id: z.uuid(),
	name: z.string(),
	position: z.number(),
	created_at: z.string(),
});
export type BoardColumn = z.infer<typeof BoardColumnSchema>;

export const BoardSchema = z.object({
	id: z.uuid(),
	name: z.string(),
	owner_user_id: z.uuid().nullable(),
	owner_team_id: z.uuid().nullable(),
	created_at: z.string(),
	updated_at: z.string(),
});
export type Board = z.infer<typeof BoardSchema>;

export const IssueSchema = z.object({
	id: z.uuid(),
	user_id: z.uuid(),
	user_name: z.string(),
	assignee_id: z.uuid().nullable(),
	assignee_name: z.string(),
	team_id: z.uuid().nullable(),
	board_column_id: z.uuid().nullable(),
	board_column_name: z.string(),
	position: z.number(),
	title: z.string().min(1, "Title is required").max(255, "Title too long"),
	description: z.string(),
	priority: PriorityType,
	due_date: z.string().nullable(),
	created_at: z.string(),
	updated_at: z.string(),
});
export type Issue = z.infer<typeof IssueSchema>;

export const CreateIssueSchema = z.object({
	title: z.string().min(1, "Title is required").max(255, "Title too long"),
	description: z.string().default(""),
	priority: PriorityType.default("Medium"),
	assignee_id: z.string().optional(),
	team_id: z.string().optional(),
	board_column_id: z.string().optional(),
	due_date: z.string().optional(),
});
export type CreateIssueInput = z.infer<typeof CreateIssueSchema>;

export const UpdateIssueSchema = z.object({
	title: z.string().min(1).max(255).optional(),
	description: z.string().optional(),
	priority: PriorityType.optional(),
	assignee_id: z.string().nullable().optional(),
	team_id: z.string().nullable().optional(),
	due_date: z.string().nullable().optional(),
});
export type UpdateIssueInput = z.infer<typeof UpdateIssueSchema>;

export const MoveIssueSchema = z.object({
	board_column_id: z.uuid(),
	position: z.number(),
});
export type MoveIssueInput = z.infer<typeof MoveIssueSchema>;
