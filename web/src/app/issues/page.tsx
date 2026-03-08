"use client";

import { usePersonalBoard } from "@/hooks/useBoards";
import { KanbanBoard } from "@/components/kanban/KanbanBoard";

export default function IssuesListPage() {
	const { data, isLoading, error } = usePersonalBoard();

	if (isLoading) {
		return (
			<div className="mx-auto max-w-6xl px-4 py-8">
				<p className="text-muted-foreground">Loading board...</p>
			</div>
		);
	}

	if (error || !data) {
		return (
			<div className="mx-auto max-w-6xl px-4 py-8">
				<p className="text-destructive">
					Failed to load board{error ? `: ${error.message}` : ""}
				</p>
			</div>
		);
	}

	return (
		<div className="mx-auto max-w-6xl px-4 py-8">
			<KanbanBoard board={data.board} columns={data.columns} />
		</div>
	);
}
