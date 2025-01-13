"use client";

import { cn } from "@repo/ui/lib/utils";
import { usePathname, useRouter } from "next/navigation";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarMenu,
  SidebarGroupContent,
  SidebarMenuButton,
  SidebarFooter,
  SidebarTrigger,
  useSidebar,
} from "@ui/components/ui/sidebar";
import SidebarAvatarMenu from "./sidebar-avatar-menu";
import { User, CircleHalf, LockSimple, ArrowUUpLeft } from "@ui/components/react-icons";
interface SidebarProps {
  className?: string;
}

export function AccountSidebar({ className }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { state } = useSidebar();
  const collapsed = state === "collapsed";

  const routes = {
    Settings: {
      children: [
        {
          name: "Profile",
          icon: <User className="mr-1 h-4 w-4" />,
          route: `/account/settings/profile`,
          key: "account-profile",
        },
        {
          name: "Security",
          icon: <LockSimple className="mr-1 h-4 w-4" />,
          route: `/account/settings/security`,
          key: "account-security",
        },
        {
          name: "Appearance",
          icon: <CircleHalf weight="fill" className="mr-1 h-4 w-4" />,
          route: `/account/settings/appearance`,
          key: "account-appearance",
        },
      ],
    },
  };

  return (
    <Sidebar collapsible="icon" className="transition-[width] duration-300 ease-in-out">
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <div className={cn("pb-12", className)}>
                <div className="space-y-4 py-4">
                  {Object.entries(routes).map(([key, route]) => (
                    <div key={key}>
                      {!collapsed && (
                        <div className="flex flex-row items-center justify-between transition-all duration-300 overflow-hidden">
                          <h2 className="mb-2 px-1 text-lg font-semibold tracking-tight">{key}</h2>
                        </div>
                      )}
                      <div className="space-y-1">
                        {route.children.map((child) => (
                          <SidebarMenuButton
                            onClick={() => router.push(child.route)}
                            key={child.key}
                            isActive={pathname === child.route}
                            tooltip={child.name}
                            className="w-full justify-start"
                          >
                            {child.icon}
                            {child.name}
                          </SidebarMenuButton>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
                <div>
                  <SidebarMenuButton
                    onClick={() => router.push("/organizations")}
                    className={`flex justify-center px-2 w-full whitespace-nowrap ${collapsed ? "mt-0" : "mt-4"}`}
                    variant="outline"
                    tooltip={"Back to Organizations"}
                  >
                    {collapsed ? <ArrowUUpLeft /> : "Back to Organizations"}
                  </SidebarMenuButton>
                </div>
              </div>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarTrigger />
        <SidebarAvatarMenu />
      </SidebarFooter>
    </Sidebar>
  );
}
