export type Brand = {
  id: number;
  name: string;
};

export type Model = {
  id: number;
  name: string;
};

export type Category = {
  id: number;
  parent?: string;
  name: string;
  path?: string;
};

export type ModelYear = {
  startyear: number;
  endyear: number;
};

export type UserInput = {
  brand?: string;
  model?: string;
  year?: string;
  category?: string;
};

export type Product = {
  id: number;
  name: string;
  category_id: number;
  description: string;
  brand: string;
  is_universal: boolean;
};
