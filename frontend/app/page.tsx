"use client";

import { DropDown } from "@/components/dropdown";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Result, tryCatch } from "@/utils/trycatch";
import { Product, UserInput } from "@/utils/types";
import axios from "axios";

import Image from "next/image";
import { useEffect, useState } from "react";

export default function Home() {
  const [userInput, setUserInput] = useState<UserInput>({
    brand: undefined,
    model: undefined,
    year: undefined,
    category: undefined,
  });

  // Fetching brands
  useEffect(() => {
    async function fetchBrands() {
      const { data, error } = await getBrands();
      if (error !== null) {
        console.error(error);
        return;
      }

      console.log(data);
    }

    async function fetchCategories() {
      const { data, error } = await getCategories();
      if (error !== null) {
        console.error(error);
        return;
      }

      console.log(data);
    }

    fetchBrands();
    fetchCategories();
  }, []);

  useEffect(() => {
    async function fetchModels() {
      if (!userInput.brand) return;
      const { data, error } = await getModelsByBrand(userInput.brand.name);

      if (error !== null) {
        console.error(error);
        return;
      }

      console.log(data);
    }

    fetchModels();

    setUserInput({
      brand: userInput.brand,
    });
  }, [userInput.brand]);

  useEffect(() => {
    async function fetchYears() {
      if (!userInput.brand || !userInput.model) return;
      const { data, error } = await getYears(userInput.brand.name, userInput.model.name);
      if (error !== null) {
        console.error(error);
        return;
      }

      console.log(data);
    }

    fetchYears();

    setUserInput({
      brand: userInput.brand,
      model: userInput.model,
    });
  }, [userInput.brand, userInput.model]);

  return (
    <>
      <div className="flex flex-col items-center justify-center w-full gap-8 py-6">
        <Card className="w-full max-w-6xl">
          <CardHeader>
            <CardTitle>Sök efter delar</CardTitle>
            <CardDescription>Sök och sortera efter delar till både märke samt specifik cross</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-5 grid-rows-1 gap-4 px-8 py-6 rounded-md bg-muted">
              <DropDown label="Märke" placeholder="Välj märke" disabled={false} input={userInput} setInput={setUserInput} type="brand" />
              <DropDown label="Modell" placeholder="Välj modell" disabled={userInput?.brand ? false : true} input={userInput} setInput={setUserInput} type="model" />
              <DropDown label="Motorstorlek" placeholder="Välj motorstorlek (cc)" disabled={userInput?.model ? false : true} input={userInput} setInput={setUserInput} type="year" />
              <DropDown label="Kategori" placeholder="Välj kategori" disabled={false} input={userInput} setInput={setUserInput} type="category" />
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
  Fetch products
*/
async function getBrands(): Promise<Result<string[], Error>> {
  const { data, error } = await tryCatch(axios.get("http://localhost:8000/brands"));
  if (error) {
    console.error("Error when fetching brands: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

async function getModelsByBrand(brand: string): Promise<Result<string[], Error>> {
  const { data, error } = await tryCatch(axios.get(`http://localhost:8000/brands/${brand}/models`));
  if (error) {
    console.error("Error when fetching models by brand: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

async function getYears(brand: string, model: string): Promise<Result<number[], Error>> {
  const { data, error } = await tryCatch(axios.get(`http://localhost:8000/brands/${brand}/models/${model}/years`));
  if (error) {
    console.error("Error when fetching years: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

export async function getCategories(): Promise<Result<string[], Error>> {
  const { data, error } = await tryCatch(axios.get("http://localhost:8000/categories"));
  if (error) {
    console.error("Error when fetching categories: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

export async function getFilteredProducts(input: UserInput): Promise<Result<Product[], Error>> {
  const params = new URLSearchParams();

  if (input.brand?.name) params.append("brand", input.brand.name);
  if (input.model?.name) params.append("model", input.model.name);
  if (input.year?.toString()) params.append("year", input.year.toString());
  if (input.category) params.append("category", input.category.name);

  const { data, error } = await tryCatch(axios.get(`http://localhost:8000/products?${params.toString()}`));

  if (error) {
    console.error("Error fetching filtered products:", error);
    return { data: null, error };
  }

  return { data: data.data, error: null };
}
