"use client";

import * as React from "react";
import Link from "next/link";

import { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList, NavigationMenuTrigger, navigationMenuTriggerStyle } from "@/components/ui/navigation-menu";
import Image from "next/image";

export function NavigationMenuDemo() {
  return (
    <NavigationMenu viewport={false} className="z-10 py-6">
      <NavigationMenuList>
        <NavigationMenuItem asChild className={navigationMenuTriggerStyle()}>
          <Image src={"/EMX.png"} width={1138} height={621} alt="Emx-Logo" />
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
            <Link href="/">Hem</Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
            <Link href="/upload">Ladda upp</Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
            <Link href="/admin">Admin</Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>
  );
}
