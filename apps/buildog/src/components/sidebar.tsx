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
  SidebarTrigger,
  SidebarFooter,
  useSidebar,
} from "@ui/components/ui/sidebar";
import SidebarAvatarMenu from "./sidebar-avatar-menu";
import { NotePencil, HashStraight, Globe, Gear } from "@ui/components/react-icons";

interface SidebarProps {
  className?: string;
  organizationId: string;
}

export function AppSidebar({ className, organizationId }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { state } = useSidebar();
  const collapsed = state === "collapsed";

  const routes = {
    Buildog: {
      children: [
        {
          name: "Blog",
          icon: <HashStraight className="mr-1 h-4 w-4" />,
          route: `/organizations/${organizationId}/blog`,
          key: "main-blog",
        },
        {
          name: "Create Blog",
          icon: <NotePencil className="mr-1 h-4 w-4" />,
          route: `/organizations/${organizationId}/create-blog`,
          key: "main-create-blog",
        },
        {
          name: "Web",
          icon: <Globe className="mr-1 h-4 w-4" />,
          route: `/organizations/${organizationId}/web`,
          key: "main-web",
        },
        {
          name: "Settings",
          icon: <Gear className="mr-1 h-4 w-4" />,
          route: `/organizations/${organizationId}/settings`,
          key: "main-settings",
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
                            className="w-full justify-start overflow-hidden whitespace-nowrap text-ellipsis transition-all duration-300"
                          >
                            {child.icon}
                            {child.name}
                          </SidebarMenuButton>
                        ))}
                      </div>
                    </div>
                  ))}
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
