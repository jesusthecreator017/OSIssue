"use client";

import { useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Trash2 } from "lucide-react";
import { toast } from "sonner";
import { useAuth } from "@/contexts/AuthContext";
import { useTeam, useDeleteTeam } from "@/hooks/useTeams";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { TeamMembers } from "@/components/teams/TeamMembers";
import { AddMemberDialog } from "@/components/teams/AddMemberDialog";

export default function TeamDetailPage() {
	const { teamId } = useParams<{ teamId: string }>();
	const { user, isLoading: authLoading } = useAuth();
	const router = useRouter();
	const { data: team, isLoading } = useTeam(teamId);
	const deleteTeam = useDeleteTeam();

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
			<div className="mx-auto max-w-4xl px-4 py-8">
				<p className="text-muted-foreground text-sm">Loading team...</p>
			</div>
		);
	}

	if (!team) {
		return (
			<div className="mx-auto max-w-4xl px-4 py-8">
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
		<div className="mx-auto max-w-4xl px-4 py-8">
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
					<Separator />

					<div className="flex items-center justify-between">
						<h2 className="text-lg font-semibold">Members</h2>
						{isCreator && <AddMemberDialog teamId={teamId} />}
					</div>

					<TeamMembers teamId={teamId} createdBy={team.created_by} />

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
