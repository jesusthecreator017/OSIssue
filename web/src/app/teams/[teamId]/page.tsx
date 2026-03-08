"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Trash2 } from "lucide-react";
import { toast } from "sonner";
import { useAuth } from "@/contexts/AuthContext";
import { useTeam, useDeleteTeam } from "@/hooks/useTeams";
import { useTeamBoard } from "@/hooks/useBoards";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { TeamMembers } from "@/components/teams/TeamMembers";
import { AddMemberDialog } from "@/components/teams/AddMemberDialog";
import { KanbanBoard } from "@/components/kanban/KanbanBoard";

export default function TeamDetailPage() {
	const { teamId } = useParams<{ teamId: string }>();
	const { user, isLoading: authLoading } = useAuth();
	const router = useRouter();
	const { data: team, isLoading } = useTeam(teamId);
	const deleteTeam = useDeleteTeam();
	const { data: boardData } = useTeamBoard(teamId);
	const [tab, setTab] = useState<"members" | "board">("members");

	useEffect(() => {
		if (!authLoading && !user) {
			router.push("/auth/login");
		}
	}, [authLoading, user, router]);

	if (authLoading || !user) {
		return null;
	}

	if (isLoading) {
		return (
			<div className="mx-auto max-w-6xl px-4 py-8">
				<p className="text-muted-foreground text-sm">Loading team...</p>
			</div>
		);
	}

	if (!team) {
		return (
			<div className="mx-auto max-w-6xl px-4 py-8">
				<p className="text-muted-foreground text-sm">Team not found.</p>
			</div>
		);
	}

	const isCreator = user.id === team.created_by;

	function handleDelete() {
		deleteTeam.mutate(teamId, {
			onSuccess: () => {
				toast.success("Team deleted");
				router.push("/teams");
			},
			onError: (err) =>
				toast.error(
					err instanceof Error ? err.message : "Failed to delete team",
				),
		});
	}

	return (
		<div className="mx-auto max-w-6xl px-4 py-8 space-y-6">
			<Card>
				<CardHeader>
					<div className="flex items-center justify-between">
						<div>
							<CardTitle className="text-2xl">{team.name}</CardTitle>
							<p className="text-muted-foreground mt-1 text-sm">
								{team.description || "No description"}
							</p>
						</div>
						<Badge variant="secondary">Max {team.max_members} members</Badge>
					</div>
				</CardHeader>
				<CardContent className="space-y-6">
					<div className="flex gap-2">
						<Button
							variant={tab === "members" ? "default" : "outline"}
							size="sm"
							onClick={() => setTab("members")}
						>
							Members
						</Button>
						<Button
							variant={tab === "board" ? "default" : "outline"}
							size="sm"
							onClick={() => setTab("board")}
						>
							Board
						</Button>
					</div>

					<Separator />

					{tab === "members" && (
						<>
							<div className="flex items-center justify-between">
								<h2 className="text-lg font-semibold">Members</h2>
								{isCreator && <AddMemberDialog teamId={teamId} />}
							</div>
							<TeamMembers teamId={teamId} createdBy={team.created_by} />
						</>
					)}

					{tab === "board" && boardData && (
						<KanbanBoard board={boardData.board} columns={boardData.columns} />
					)}

					{tab === "board" && !boardData && (
						<p className="text-muted-foreground text-sm">Loading board...</p>
					)}

					{isCreator && (
						<>
							<Separator />
							<div>
								<h2 className="mb-2 text-lg font-semibold text-destructive">
									Danger Zone
								</h2>
								<Button
									variant="destructive"
									onClick={handleDelete}
									disabled={deleteTeam.isPending}
								>
									<Trash2 className="mr-2 size-4" />
									{deleteTeam.isPending ? "Deleting..." : "Delete Team"}
								</Button>
							</div>
						</>
					)}
				</CardContent>
			</Card>
		</div>
	);
}
