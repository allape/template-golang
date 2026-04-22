import { BaseSearchParams } from "@allape/gocrud";
import { CrudySelector } from "@allape/gocrud-react";
import { ICrudySelectorProps } from "@allape/gocrud-react/src/component/CrudySelector";
import { PropsWithChildren, ReactElement, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { TagCrudy } from "../../api/tag.ts";
import { ITag, ITagSearchParams } from "../../model/tag.ts";

type IRecord = ITag;
type ISearchParams = ITagSearchParams;

export type ITagSelectorSelectorProps = Partial<ICrudySelectorProps<IRecord>>;

export default function TagSelector(
  props: PropsWithChildren<ITagSelectorSelectorProps>,
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
      placeholder={`${t("select")} ${t("tag._")}`}
      {...props}
      crudy={TagCrudy}
      searchParams={sp}
    />
  );
}
