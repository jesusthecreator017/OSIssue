"use client";

import { useState } from "react";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { CreateIssueSchema } from "@/schemas/issue";
import { useCreateIssue } from "@/hooks/useIssues";

interface CreateIssueDialogProps {
	open: boolean;
	onOpenChange: (open: boolean) => void;
	defaultColumnId?: string;
}

export function CreateIssueDialog({
	open,
	onOpenChange,
	defaultColumnId,
}: CreateIssueDialogProps) {
	const [title, setTitle] = useState("");
	const [description, setDescription] = useState("");
	const [priority, setPriority] = useState("Medium");
	const [errors, setErrors] = useState<Record<string, string>>({});
	const createIssue = useCreateIssue();

	function resetForm() {
		setTitle("");
		setDescription("");
		setPriority("Medium");
		setErrors({});
	}

	function handleSubmit(e: React.FormEvent) {
		e.preventDefault();

		const result = CreateIssueSchema.safeParse({
			title,
			description,
			priority,
			board_column_id: defaultColumnId,
		});
		if (!result.success) {
			const fieldErrors: Record<string, string> = {};
			for (const issue of result.error.issues) {
				const key = String(issue.path[0]);
				if (key && !fieldErrors[key]) {
					fieldErrors[key] = issue.message;
				}
			}
			setErrors(fieldErrors);
			return;
		}

		setErrors({});
		createIssue.mutate(result.data, {
			onSuccess: () => {
				resetForm();
				onOpenChange(false);
			},
		});
	}

	return (
		<Dialog
			open={open}
			onOpenChange={(v) => {
				if (!v) resetForm();
				onOpenChange(v);
			}}
		>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Create Issue</DialogTitle>
					<DialogDescription>Add a new issue to the board.</DialogDescription>
				</DialogHeader>
				<form onSubmit={handleSubmit} className="grid gap-4">
					<div className="grid gap-2">
						<Label htmlFor="title">Title</Label>
						<Input
							id="title"
							value={title}
							onChange={(e) => setTitle(e.target.value)}
							placeholder="Issue title"
						/>
						{errors.title && (
							<p className="text-destructive text-sm">{errors.title}</p>
						)}
					</div>
					<div className="grid gap-2">
						<Label htmlFor="description">Description</Label>
						<Textarea
							id="description"
							value={description}
							onChange={(e) => setDescription(e.target.value)}
							placeholder="Optional description"
							rows={3}
						/>
					</div>
					<div className="grid gap-2">
						<Label htmlFor="priority">Priority</Label>
						<Select value={priority} onValueChange={setPriority}>
							<SelectTrigger>
								<SelectValue />
							</SelectTrigger>
							<SelectContent>
								<SelectItem value="Low">Low</SelectItem>
								<SelectItem value="Medium">Medium</SelectItem>
								<SelectItem value="High">High</SelectItem>
								<SelectItem value="Critical">Critical</SelectItem>
							</SelectContent>
						</Select>
					</div>
					<DialogFooter>
						<Button type="submit" disabled={createIssue.isPending}>
							{createIssue.isPending ? "Creating..." : "Create Issue"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}
