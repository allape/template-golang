import { BaseSearchParams } from "@allape/gocrud";
import { ICrudySelectorProps, PagedCrudySelector } from "@allape/gocrud-react";
import { PropsWithChildren, ReactElement, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { TagCrudy } from "../../api/tag.ts";
import { ITag, ITagSearchParams } from "../../model/tag.ts";

type IRecord = ITag;
type ISearchParams = ITagSearchParams;

export type ITagSelectorSelectorProps = Partial<ICrudySelectorProps<IRecord>>;

export default function PagedTagSelector(
  props: PropsWithChildren<ITagSelectorSelectorProps>,
): ReactElement {
  const { t } = useTranslation();

  const sp = useMemo<ISearchParams>(
    () => ({
      ...BaseSearchParams,
      orderBy_priority: "desc",
    }),
    [],
  );

  return (
    <PagedCrudySelector<IRecord, ISearchParams>
      placeholder={`${t("select")} ${t("tag._")}`}
      {...props}
      crudy={TagCrudy}
      pageSize={1000}
      searchParams={sp}
      searchPropName="like_name"
      inKeyword="in_id"
    />
  );
}
