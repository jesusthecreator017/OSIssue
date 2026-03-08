"use client";

import { useParams } from "next/navigation";
import { useIssue } from "@/hooks/useIssues";
import { useIssueLabels } from "@/hooks/useLabels";
import { Badge } from "@/components/ui/badge";
import { Calendar } from "lucide-react";
import type { PriorityType } from "@/schemas/issue";

const priorityColors: Record<PriorityType, string> = {
	Low: "bg-slate-100 text-slate-700 border-slate-200",
	Medium: "bg-blue-100 text-blue-700 border-blue-200",
	High: "bg-orange-100 text-orange-700 border-orange-200",
	Critical: "bg-red-100 text-red-700 border-red-200",
};

export default function IssuePage() {
	const { issueId } = useParams<{ issueId: string }>();
	const { data: issue, isLoading, error } = useIssue(issueId);
	const { data: labels } = useIssueLabels(issueId);

	if (isLoading) {
		return <div className="mx-auto max-w-2xl px-4 py-8">Loading...</div>;
	}

	if (error || !issue) {
		return (
			<div className="mx-auto max-w-2xl px-4 py-8">
				<p className="text-destructive">
					Failed to load issue{error ? `: ${error.message}` : ""}
				</p>
			</div>
		);
	}

	const isOverdue = issue.due_date && new Date(issue.due_date) < new Date();

	return (
		<div className="mx-auto max-w-2xl px-4 py-8 space-y-6">
			<div>
				<h1 className="text-2xl font-bold">{issue.title}</h1>
				<p className="mt-2 text-muted-foreground">
					{issue.description || "No description"}
				</p>
			</div>

			<div className="flex flex-wrap gap-2">
				<Badge className={priorityColors[issue.priority]} variant="outline">
					{issue.priority}
				</Badge>
				{issue.board_column_name && (
					<Badge variant="secondary">{issue.board_column_name}</Badge>
				)}
				{issue.due_date && (
					<Badge
						variant="outline"
						className={isOverdue ? "text-red-600 border-red-200" : ""}
					>
						<Calendar className="size-3 mr-1" />
						{new Date(issue.due_date).toLocaleDateString()}
					</Badge>
				)}
			</div>

			{labels && labels.length > 0 && (
				<div className="flex flex-wrap gap-1.5">
					{labels.map((label) => (
						<Badge
							key={label.id}
							variant="outline"
							style={{ borderColor: label.color, color: label.color }}
						>
							{label.name}
						</Badge>
					))}
				</div>
			)}

			<div className="grid grid-cols-2 gap-4 text-sm">
				<div>
					<span className="text-muted-foreground">Created by: </span>
					<span>{issue.user_name}</span>
				</div>
				<div>
					<span className="text-muted-foreground">Assignee: </span>
					<span>{issue.assignee_name || "Unassigned"}</span>
				</div>
				<div>
					<span className="text-muted-foreground">Created: </span>
					<span>{new Date(issue.created_at).toLocaleDateString()}</span>
				</div>
				<div>
					<span className="text-muted-foreground">Updated: </span>
					<span>{new Date(issue.updated_at).toLocaleDateString()}</span>
				</div>
			</div>
		</div>
	);
}
