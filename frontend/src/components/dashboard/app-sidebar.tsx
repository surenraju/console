"use client"

import * as React from "react"
import {
  IconDashboard,
  IconRoute,
  IconServer,
  IconShield,
  IconSettings,
  IconRobot,
} from "@tabler/icons-react"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { NavMain } from "./nav-main"
import { NavSecondary } from "./nav-secondary"
import { NavUser } from "./nav-user"

const data = {
  user: {
    name: "Admin",
    email: "admin@envoy-ai-gateway.com",
    avatar: "/avatars/admin.jpg",
  },
  navMain: [
    {
      title: "Dashboard",
      path: "/",
      icon: IconDashboard,
    },
    {
      title: "Routes",
      path: "/routes",
      icon: IconRoute,
    },
    {
      title: "Backends",
      path: "/backends",
      icon: IconServer,
    },
    {
      title: "LLM Providers",
      path: "/llm-providers",
      icon: IconRobot,
    },
    {
      title: "Policies",
      path: "/policies",
      icon: IconShield,
    },
  ],
  navSecondary: [
    {
      title: "Settings",
      path: "/settings",
      icon: IconSettings,
    },
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className="data-[slot=sidebar-menu-button]:!p-1.5"
            >
              <a href="#">
                <IconRoute className="!size-5" />
                <span className="text-base font-semibold">Envoy AI Gateway</span>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
    </Sidebar>
  )
} 