"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
	Bug,
	Menu,
	LayoutDashboard,
	AlertCircle,
	Shield,
	LogOut,
	User,
	Users,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
	Sheet,
	SheetContent,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet";
import { Separator } from "@/components/ui/separator";
import { UserAvatar } from "@/components/UserAvatar";
import { useAuth } from "@/contexts/AuthContext";
import { hasPermission, PERM_ADMIN } from "@/schemas/user";

export function Navbar() {
	const { user, isLoading, logout } = useAuth();
	const router = useRouter();
	const [open, setOpen] = useState(false);

	// Nav bar links
	const navLinks = [
		{ name: "Dashboard", href: "/", icon: <LayoutDashboard />},
		{ name: "Profile", href: "/profile", icon: <User />},
		{ name: "Issues", href: "/issues", icon: <AlertCircle />},
		{ name: "Teams", href: "/teams", icon: <Users />},
	];
	const adminLinks = [
		{ name: "Stats", href: "/admin", icon: <Shield />},
	]

	function handleLogout() {
		logout();
		setOpen(false);
		router.push("/");
	}

	if (isLoading) {
		return (
			<nav className="border-b">
				<div className="mx-auto flex h-14 max-w-6xl items-center px-4">
					<Bug className="size-6" />
				</div>
			</nav>
		);
	}

	const isAdmin = user ? hasPermission(user.permissions, PERM_ADMIN) : false;

	return (
		<nav className="border-b">
			<div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
				<Link href="/" className="flex items-center gap-2">
					<Bug className="size-6" />
				</Link>

				{!user ? (
					<div className="flex items-center gap-2">
						<Button variant="ghost" asChild>
							<Link href="/auth/login">Login</Link>
						</Button>
						<Button asChild>
							<Link href="/auth/register">Register</Link>
						</Button>
					</div>
				) : (
					<Sheet open={open} onOpenChange={setOpen}>
						<SheetTrigger asChild>
							<Button variant="ghost" size="icon">
								<Menu className="size-5" />
							</Button>
						</SheetTrigger>
						<SheetContent side="right" className="flex w-72 flex-col">
							<SheetTitle className="sr-only">Navigation menu</SheetTitle>

							<div className="flex items-center gap-3 py-4 px-3">
								<UserAvatar name={user.name} />
								<div className="flex flex-col">
									<span className="text-sm font-medium">{user.name}</span>
									<span className="text-muted-foreground text-xs">
										{user.email}
									</span>
								</div>
							</div>

							<Separator />

							<div className="flex flex-col gap-1 py-2 px-3">
								<span className="text-muted-foreground px-1 text-xs font-semibold uppercase">
									Navigation
								</span>
								{navLinks.map( tab => (
									<Button
									key={tab.name}
									variant="ghost"
									className="justify-start"
									asChild
									onClick={() => setOpen(false)}
								>
									<Link href={tab.href}>
										{tab.icon}
										{tab.name}
									</Link>
								</Button>
								))}
							</div>

							{isAdmin && (
								<>
									<Separator />
									<div className="flex flex-col gap-1 py-2 px-3">
										<span className="text-muted-foreground px-1 text-xs font-semibold uppercase">
											Admin
										</span>
										{adminLinks.map( link => (
											<Button
											key={link.name}
											variant="ghost"
											className="justify-start"
											asChild
											onClick={() => setOpen(false)}
										>
											<Link href={link.href}>
												{link.icon}
												{link.name}
											</Link>
										</Button>
										))}
									</div>
								</>
							)}

							<div className="mt-auto">
								<Separator />
								<Button
									variant="ghost"
									className="w-full justify-start py-4"
									onClick={handleLogout}
								>
									<LogOut className="mr-2 size-4" />
									Logout
								</Button>
							</div>
						</SheetContent>
					</Sheet>
				)}
			</div>
		</nav>
	);
}
