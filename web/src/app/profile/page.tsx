"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Mail, Calendar, Shield, User as UserIcon } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import {
	hasPermission,
	PERM_READ,
	PERM_WRITE,
	PERM_ADMIN,
} from "@/schemas/user";
import { UserAvatar } from "@/components/UserAvatar";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

function formatDate(dateStr: string) {
	return new Date(dateStr).toLocaleDateString("en-US", {
		year: "numeric",
		month: "long",
		day: "numeric",
	});
}

export default function ProfilePage() {
	const { user, isLoading: authLoading } = useAuth();
	const router = useRouter();

	useEffect(() => {
		if (!authLoading && !user) {
			router.push("/auth/login");
		}
	}, [authLoading, user, router]);

	if (authLoading || !user) {
		return null;
	}

	const permissions: string[] = [];
	if (hasPermission(user.permissions, PERM_READ)) permissions.push("Read");
	if (hasPermission(user.permissions, PERM_WRITE)) permissions.push("Write");
	if (hasPermission(user.permissions, PERM_ADMIN)) permissions.push("Admin");

	return (
		<div className="mx-auto max-w-2xl px-4 py-8">
			<h1 className="mb-6 text-2xl font-bold">Profile</h1>

			<Card>
				<CardHeader>
					<div className="flex items-center gap-4">
						<UserAvatar name={user.name} className="size-16 text-xl" />
						<div>
							<CardTitle className="text-xl">{user.name}</CardTitle>
							<p className="text-muted-foreground text-sm">{user.email}</p>
						</div>
					</div>
				</CardHeader>
				<CardContent className="space-y-4">
					<Separator />

					<div className="grid gap-4 sm:grid-cols-2">
						<div className="flex items-center gap-2">
							<UserIcon className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">Name</span>
						</div>
						<span className="text-sm font-medium">{user.name}</span>

						<div className="flex items-center gap-2">
							<Mail className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">Email</span>
						</div>
						<span className="text-sm font-medium">{user.email}</span>

						<div className="flex items-center gap-2">
							<Shield className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">Permissions</span>
						</div>
						<div className="flex gap-1">
							{permissions.map((perm) => (
								<Badge key={perm} variant="secondary">
									{perm}
								</Badge>
							))}
						</div>

						<div className="flex items-center gap-2">
							<Calendar className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">Joined</span>
						</div>
						<span className="text-sm font-medium">
							{formatDate(user.created_at)}
						</span>

						<div className="flex items-center gap-2">
							<Calendar className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">
								Last updated
							</span>
						</div>
						<span className="text-sm font-medium">
							{formatDate(user.updated_at)}
						</span>
					</div>
				</CardContent>
			</Card>
		</div>
	);
}
