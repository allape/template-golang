import Crudy, { config } from "@allape/gocrud-react";
import { IGallery } from "../model/gallery.ts";

export const GalleryCrudy = new Crudy<IGallery>(`${config.SERVER_URL}/gallery`);
