"use client";

import { useDroppable } from "@dnd-kit/core";
import {
	SortableContext,
	verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { Badge } from "@/components/ui/badge";
import { IssueCard } from "./IssueCard";
import type { BoardColumn, Issue } from "@/schemas/issue";

interface KanbanColumnProps {
	column: BoardColumn;
	issues: Issue[];
	boardId: string;
}

export function KanbanColumn({ column, issues, boardId: _boardId }: KanbanColumnProps) {
	const { setNodeRef, isOver } = useDroppable({ id: column.id });

	return (
		<div
			ref={setNodeRef}
			className={`flex flex-1 flex-col gap-3 min-w-70 rounded-lg p-2 transition-colors ${
				isOver ? "bg-accent/50" : ""
			}`}
		>
			<div className="flex items-center gap-2 px-1">
				<h2 className="text-sm font-semibold">{column.name}</h2>
				<Badge variant="secondary">{issues.length}</Badge>
			</div>
			<SortableContext
				items={issues.map((i) => i.id)}
				strategy={verticalListSortingStrategy}
			>
				<div className="flex flex-col gap-2">
					{issues.map((issue) => (
						<IssueCard key={issue.id} issue={issue} />
					))}
					{issues.length === 0 && (
						<p className="text-muted-foreground text-center text-sm py-8">
							No issues
						</p>
					)}
				</div>
			</SortableContext>
		</div>
	);
}
