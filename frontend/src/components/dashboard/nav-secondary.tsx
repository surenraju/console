"use client"

import * as React from "react"
import { cn } from "@/lib/utils"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { Link, useLocation } from "react-router-dom"

interface NavSecondaryProps {
  items: {
    title: string
    path: string
    icon: React.ComponentType<{ className?: string }>
  }[]
  className?: string
}

export function NavSecondary({ items, className }: NavSecondaryProps) {
  const location = useLocation()

  return (
    <SidebarMenu className={cn("mt-auto", className)}>
      {items.map((item) => (
        <SidebarMenuItem key={item.title}>
          <SidebarMenuButton 
            asChild
            isActive={location.pathname === item.path}
          >
            <Link to={item.path}>
              <item.icon className="!size-4" />
              <span>{item.title}</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      ))}
    </SidebarMenu>
  )
} 