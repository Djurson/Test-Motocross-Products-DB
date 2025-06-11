"use client";

import { DropDown } from "@/components/dropdown";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { tryCatch } from "@/utils/trycatch";
import { Product, UserInput } from "@/utils/types";
import axios from "axios";

import Image from "next/image";
import { useEffect, useState } from "react";

export default function Home() {
  const [userInput, setUserInput] = useState<UserInput>({
    brand: undefined,
    model: undefined,
    engineSize: undefined,
    year: undefined,
    category: undefined,
  });

  /*
    Implementera funktioner för att hämta märken, motorstorlekar osv osv
  */

  return (
    <>
      <div className="flex flex-col items-center justify-center w-full gap-8 py-6">
        <Image src={"/EMX.png"} alt="EMX logo" width={1138} height={621} priority className="w-1/8" />
        <Card className="w-full max-w-6xl">
          <CardHeader>
            <CardTitle>Sök efter delar</CardTitle>
            <CardDescription>Sök och sortera efter delar till både märke samt specifik cross</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-5 grid-rows-1 gap-4 px-8 py-6 rounded-md bg-muted">
              <DropDown label="Märke" placeholder="Välj märke" disabled={false} input={userInput} setInput={setUserInput} type="brand" />
              <DropDown label="Modell" placeholder="Välj modell" disabled={userInput?.brand ? false : true} input={userInput} setInput={setUserInput} type="model" />
              <DropDown label="Motorstorlek" placeholder="Välj motorstorlek (cc)" disabled={userInput?.model ? false : true} input={userInput} setInput={setUserInput} type="engineSize" />
              <DropDown label="Kategori" placeholder="Välj kategori" disabled={userInput?.engineSize ? false : true} input={userInput} setInput={setUserInput} type="category" />
              <Button className="self-end" variant="default">
                Sök
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </>
  );
}

/* 
  Omvandla till att använda tryCatch samt axios 
*/
export async function getProductsByBrand(brand: string): Promise<Product[]> {
  const res = await fetch(`http://localhost:8080/brands/${brand}/products`);
  if (!res.ok) throw new Error("Failed to fetch products by brand");
  return res.json();
}

export async function getEngineSizes(brand: string, model: string): Promise<string[]> {
  const res = await fetch(`http://localhost:8080/brands/${brand}/models/${model}/engine-sizes`);
  if (!res.ok) throw new Error("Failed to fetch engine sizes");
  return res.json();
}

export async function getProductsByBrandModel(brand: string, model: string): Promise<Product[]> {
  const res = await fetch(`http://localhost:8080/brands/${brand}/models/${model}/products`);
  if (!res.ok) throw new Error("Failed to fetch products by brand and model");
  return res.json();
}

export async function getYears(brand: string, model: string, engineSize: string): Promise<number[]> {
  const res = await fetch(`http://localhost:8080/brands/${brand}/models/${model}/engine-sizes/${engineSize}/years`);
  if (!res.ok) throw new Error("Failed to fetch years");
  return res.json();
}

export async function getProductsByBrandModelEngine(brand: string, model: string, engineSize: string): Promise<Product[]> {
  const res = await fetch(`http://localhost:8080/brands/${brand}/models/${model}/engine-sizes/${engineSize}/products`);
  if (!res.ok) throw new Error("Failed to fetch products by brand, model, and engine");
  return res.json();
}

export async function getProductsFullFilter(brand: string, model: string, engineSize: string, year: number): Promise<Product[]> {
  const res = await fetch(`http://localhost:8080/brands/${brand}/models/${model}/engine-sizes/${engineSize}/years/${year}/products`);
  if (!res.ok) throw new Error("Failed to fetch full-filtered products");
  return res.json();
}

export async function getBrands(): Promise<string[]> {
  const res = await fetch("http://localhost:8080/brands");
  if (!res.ok) throw new Error("Failed to fetch brands");
  return res.json();
}

export async function getModelsByBrand(brand: string): Promise<string[]> {
  const res = await fetch(`http://localhost:8080/brands/${brand}/models`);
  if (!res.ok) throw new Error("Failed to fetch models");
  return res.json();
}
