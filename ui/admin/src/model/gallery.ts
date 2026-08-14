import { IBaseSearchParams } from "@allape/gocrud";

export interface IGallerySearchParams extends IBaseSearchParams {
  like_name?: string;
  like_keywords?: string;
  keywords?: string;
  isPublic?: boolean;
  createdBy?: string;
}
