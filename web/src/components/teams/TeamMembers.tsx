"use client";

import { Trash2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { UserAvatar } from "@/components/UserAvatar";
import { useTeamMembers, useRemoveTeamMember } from "@/hooks/useTeams";
import { useAuth } from "@/contexts/AuthContext";
import type { TeamMember } from "@/schemas/teams";

const roleBadgeVariant: Record<string, "default" | "secondary" | "destructive"> = {
	owner: "default",
	admin: "secondary",
	member: "secondary",
};

export function TeamMembers({
	teamId,
	createdBy,
}: {
	teamId: string;
	createdBy: string;
}) {
	const { user } = useAuth();
	const { data: members, isLoading } = useTeamMembers(teamId);
	const removeMember = useRemoveTeamMember(teamId);

	const isCreator = user?.id === createdBy;

	function handleRemove(member: TeamMember) {
		removeMember.mutate(member.user_id, {
			onSuccess: () => toast.success("Member removed"),
			onError: (err) =>
				toast.error(err instanceof Error ? err.message : "Failed to remove member"),
		});
	}

	if (isLoading) {
		return <p className="text-muted-foreground text-sm">Loading members...</p>;
	}

	if (!members?.length) {
		return <p className="text-muted-foreground text-sm">No members yet.</p>;
	}

	return (
		<div className="space-y-2">
			{members.map((member) => (
				<div
					key={member.user_id}
					className="flex items-center justify-between rounded-md border p-3"
				>
					<div className="flex items-center gap-3">
						<UserAvatar name={member.user_name ?? member.user_id.slice(0, 8)} />
						<div>
							<p className="text-sm font-medium">{member.user_name ?? member.user_id}</p>
							{member.email && <p className="text-muted-foreground text-xs">{member.email}</p>}
							<p className="text-muted-foreground text-xs">
								Joined{" "}
								{new Date(member.joined_at).toLocaleDateString()}
							</p>
						</div>
					</div>
					<div className="flex items-center gap-2">
						<Badge variant={roleBadgeVariant[member.role] ?? "secondary"}>
							{member.role}
						</Badge>
						{isCreator && member.user_id !== user?.id && (
							<Button
								variant="ghost"
								size="icon"
								onClick={() => handleRemove(member)}
								disabled={removeMember.isPending}
							>
								<Trash2 className="size-4" />
							</Button>
						)}
					</div>
				</div>
			))}
		</div>
	);
}
