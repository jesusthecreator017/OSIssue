"use client";

import { useState, useMemo } from "react";
import {
	DndContext,
	DragOverlay,
	closestCorners,
	PointerSensor,
	useSensor,
	useSensors,
	type DragStartEvent,
	type DragEndEvent,
	type DragOverEvent,
} from "@dnd-kit/core";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useBoardIssues } from "@/hooks/useBoards";
import { useMoveIssue } from "@/hooks/useIssues";
import { KanbanColumn } from "./KanbanColumn";
import { IssueCard } from "./IssueCard";
import { CreateIssueDialog } from "@/components/issues/CreateIssueDialog";
import type { Board, BoardColumn, Issue } from "@/schemas/issue";

interface KanbanBoardProps {
	board: Board;
	columns: BoardColumn[];
}

export function KanbanBoard({ board, columns }: KanbanBoardProps) {
	const [createOpen, setCreateOpen] = useState(false);
	const [activeIssue, setActiveIssue] = useState<Issue | null>(null);
	const {
		data: issues,
		isLoading,
		error,
	} = useBoardIssues(board.id);
	const moveIssue = useMoveIssue();

	const sensors = useSensors(
		useSensor(PointerSensor, {
			activationConstraint: { distance: 8 },
		}),
	);

	const issuesByColumn = useMemo(() => {
		const map: Record<string, Issue[]> = {};
		for (const col of columns) {
			map[col.id] = [];
		}
		for (const issue of issues ?? []) {
			const colId = issue.board_column_id;
			if (colId && map[colId]) {
				map[colId].push(issue);
			}
		}
		// Sort each column's issues by position
		for (const colId of Object.keys(map)) {
			map[colId].sort((a, b) => a.position - b.position);
		}
		return map;
	}, [issues, columns]);

	function handleDragStart(event: DragStartEvent) {
		const issue = (issues ?? []).find((i) => i.id === event.active.id);
		setActiveIssue(issue ?? null);
	}

	function handleDragOver(_event: DragOverEvent) {
		// Could add optimistic reordering here
	}

	function handleDragEnd(event: DragEndEvent) {
		setActiveIssue(null);
		const { active, over } = event;
		if (!over) return;

		const issueId = active.id as string;
		let targetColumnId: string;
		let targetPosition: number;

		// Determine if we dropped on a column or on another issue
		const overColumnId = columns.find((c) => c.id === over.id)?.id;
		if (overColumnId) {
			// Dropped directly on a column
			targetColumnId = overColumnId;
			targetPosition = (issuesByColumn[targetColumnId]?.length ?? 0);
		} else {
			// Dropped on an issue — find its column
			const overIssue = (issues ?? []).find((i) => i.id === over.id);
			if (!overIssue?.board_column_id) return;
			targetColumnId = overIssue.board_column_id;
			targetPosition = overIssue.position;
		}

		const currentIssue = (issues ?? []).find((i) => i.id === issueId);
		if (
			currentIssue?.board_column_id === targetColumnId &&
			currentIssue?.position === targetPosition
		) {
			return;
		}

		moveIssue.mutate({
			id: issueId,
			data: {
				board_column_id: targetColumnId,
				position: targetPosition,
			},
		});
	}

	if (isLoading) {
		return (
			<div className="flex items-center justify-center py-20">
				<p className="text-muted-foreground">Loading issues...</p>
			</div>
		);
	}

	if (error) {
		return (
			<div className="flex items-center justify-center py-20">
				<p className="text-destructive">
					Failed to load issues: {error.message}
				</p>
			</div>
		);
	}

	return (
		<div className="grid gap-6">
			<div className="flex items-center justify-between">
				<h1 className="text-2xl font-bold">Issues</h1>
				<Button onClick={() => setCreateOpen(true)}>
					<Plus className="size-4 mr-1" />
					New Issue
				</Button>
			</div>
			<DndContext
				sensors={sensors}
				collisionDetection={closestCorners}
				onDragStart={handleDragStart}
				onDragOver={handleDragOver}
				onDragEnd={handleDragEnd}
			>
				<div className="flex gap-4 overflow-x-auto pb-4">
					{columns.map((col) => (
						<KanbanColumn
							key={col.id}
							column={col}
							issues={issuesByColumn[col.id] ?? []}
							boardId={board.id}
						/>
					))}
				</div>
				<DragOverlay>
					{activeIssue ? <IssueCard issue={activeIssue} isDragOverlay /> : null}
				</DragOverlay>
			</DndContext>
			<CreateIssueDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				defaultColumnId={columns[0]?.id}
			/>
		</div>
	);
}
