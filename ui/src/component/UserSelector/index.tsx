import { BaseSearchParams } from "@allape/gocrud";
import { CrudySelector } from "@allape/gocrud-react";
import { ICrudySelectorProps } from "@allape/gocrud-react/src/component/CrudySelector";
import { PropsWithChildren, ReactElement, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { UserCrudy } from "../../api/user.ts";
import { IUser, IUserSearchParams } from "../../model/user.ts";

type IRecord = IUser;
type ISearchParams = IUserSearchParams;

export type IUserSelectorSelectorProps = Partial<ICrudySelectorProps<IRecord>>;

export default function UserSelector(
  props: PropsWithChildren<IUserSelectorSelectorProps>,
): ReactElement {
  const { t } = useTranslation();

  const sp = useMemo<ISearchParams>(
    () => ({
      ...BaseSearchParams,
      orderBy_index: "asc",
    }),
    [],
  );

  return (
    <CrudySelector<IRecord, ISearchParams>
      placeholder={`${t("select")} ${t("user._")}`}
      {...props}
      crudy={UserCrudy}
      pageSize={1000}
      searchParams={sp}
      searchPropName="name"
    />
  );
}
