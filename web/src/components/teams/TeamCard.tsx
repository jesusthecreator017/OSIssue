"use client";

import Link from "next/link";
import { Users } from "lucide-react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { Team } from "@/schemas/teams";

function formatDate(dateStr: string) {
	return new Date(dateStr).toLocaleDateString("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric",
	});
}

export function TeamCard({ team }: { team: Team }) {
	return (
		<Link href={`/teams/${team.id}`}>
			<Card className="transition-colors hover:border-foreground/20">
				<CardHeader className="pb-2">
					<div className="flex items-center justify-between">
						<CardTitle className="text-lg">{team.name}</CardTitle>
						<Badge variant="secondary">
							<Users className="mr-1 size-3" />
							{team.max_members}
						</Badge>
					</div>
				</CardHeader>
				<CardContent>
					<p className="text-muted-foreground line-clamp-2 text-sm">
						{team.description || "No description"}
					</p>
					<p className="text-muted-foreground mt-2 text-xs">
						Created {formatDate(team.created_at)}
					</p>
				</CardContent>
			</Card>
		</Link>
	);
}
