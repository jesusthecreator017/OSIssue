"use client";

import { useState } from "react";
import Link from "next/link";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Trash2, Calendar } from "lucide-react";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
	CardAction,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DeleteIssueAlert } from "@/components/issues/DeleteIssueAlert";
import type { Issue, PriorityType } from "@/schemas/issue";

const priorityColors: Record<PriorityType, string> = {
	Low: "bg-slate-100 text-slate-700 border-slate-200",
	Medium: "bg-blue-100 text-blue-700 border-blue-200",
	High: "bg-orange-100 text-orange-700 border-orange-200",
	Critical: "bg-red-100 text-red-700 border-red-200",
};

interface IssueCardProps {
	issue: Issue;
	isDragOverlay?: boolean;
}

export function IssueCard({ issue, isDragOverlay }: IssueCardProps) {
	const [deleteOpen, setDeleteOpen] = useState(false);
	const {
		attributes,
		listeners,
		setNodeRef,
		transform,
		transition,
		isDragging,
	} = useSortable({ id: issue.id });

	const style = {
		transform: CSS.Transform.toString(transform),
		transition,
		opacity: isDragging ? 0.5 : 1,
	};

	const isOverdue =
		issue.due_date && new Date(issue.due_date) < new Date();

	return (
		<>
			<div
				ref={setNodeRef}
				style={isDragOverlay ? undefined : style}
				{...attributes}
				{...listeners}
			>
				<Card className={`gap-3 py-4 cursor-grab active:cursor-grabbing ${isDragOverlay ? "shadow-lg rotate-2" : ""}`}>
					<CardHeader className="gap-1">
						<CardTitle className="text-sm">
							<Link
								href={`/issues/${issue.id}`}
								className="text-foreground no-underline hover:underline"
							>
								{issue.title}
							</Link>
						</CardTitle>
						<CardAction>
							<Button
								variant="ghost"
								size="icon-xs"
								onClick={(e) => {
									e.stopPropagation();
									setDeleteOpen(true);
								}}
							>
								<Trash2 className="size-3.5 text-muted-foreground" />
							</Button>
						</CardAction>
					</CardHeader>
					<CardContent className="grid gap-2">
						{issue.description && (
							<p className="text-muted-foreground text-xs line-clamp-2">
								{issue.description}
							</p>
						)}
						<div className="flex items-center flex-wrap gap-1.5">
							<Badge
								className={priorityColors[issue.priority]}
								variant="outline"
							>
								{issue.priority}
							</Badge>
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
						<div className="flex items-center justify-between gap-2">
							{issue.assignee_name ? (
								<span className="text-xs text-muted-foreground truncate">
									{issue.assignee_name}
								</span>
							) : (
								<span className="text-xs text-muted-foreground/50 truncate">
									Unassigned
								</span>
							)}
							<span className="text-muted-foreground text-xs truncate">
								{issue.user_name}
							</span>
						</div>
					</CardContent>
				</Card>
			</div>
			<DeleteIssueAlert
				issueId={issue.id}
				issueTitle={issue.title}
				open={deleteOpen}
				onOpenChange={setDeleteOpen}
			/>
		</>
	);
}
