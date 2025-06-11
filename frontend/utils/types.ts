export type Brand = {
  id: number;
  name: string;
};

export type Model = {
  id: number;
  name: string;
};

export type EngineSize = {
  id: number;
  sizeCC: number;
};

export type Category = {
  id: number;
  parent?: string;
  category: string;
};

export type ModelYear = {
  startYear: number;
  endYear: number;
};

export type UserInput = {
  brand: Brand | undefined;
  model?: Model;
  engineSize?: EngineSize;
  year?: ModelYear;
  category?: Category;
};

export type Product = {
  id: number;
  name: string;
  category_id: number;
  description: string;
  brand: string;
  is_universal: boolean;
};
