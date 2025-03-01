import Crudy, { config } from "@allape/gocrud-react";
import { IUser } from "../model/user.ts";

export const UserCrudy = new Crudy<IUser>(`${config.SERVER_URL}/user`);
