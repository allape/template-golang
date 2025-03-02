import { BaseSearchParams } from "@allape/gocrud";
import { ICrudySelectorProps, PagedCrudySelector } from "@allape/gocrud-react";
import { PropsWithChildren, ReactElement, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { UserCrudy } from "../../api/user.ts";
import { IUser, IUserSearchParams } from "../../model/user.ts";

type IRecord = IUser;
type ISearchParams = IUserSearchParams;

export type IUserSelectorSelectorProps = Partial<ICrudySelectorProps<IRecord>>;

export default function PagedUserSelector(
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
    <PagedCrudySelector<IRecord, ISearchParams>
      placeholder={`${t("select")} ${t("user._")}`}
      {...props}
      crudy={UserCrudy}
      pageSize={1000}
      searchParams={sp}
      searchPropName="like_name"
      inKeyword="in_id"
    />
  );
}
