import { BaseSearchParams } from "@allape/gocrud";
import { ICrudySelectorProps, PagedCrudySelector } from "@allape/gocrud-react";
import { PropsWithChildren, ReactElement, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { GalleryCrudy } from "../../api/gallery.ts";
import { IGallery, IGallerySearchParams } from "../../model/gallery.ts";

type IRecord = IGallery;
type ISearchParams = IGallerySearchParams;

export type IGallerySelectorProps = Partial<
  ICrudySelectorProps<IRecord>
>;

export default function GallerySelector(
  props: PropsWithChildren<IGallerySelectorProps>,
): ReactElement {
  const { t } = useTranslation();

  const sp = useMemo<ISearchParams>(
    () => ({
      ...BaseSearchParams,
      orderByDefault: "1",
    }),
    [],
  );

  return (
    <PagedCrudySelector<IRecord, ISearchParams>
      placeholder={`${t("select")} ${t("gallery._")}`}
      {...props}
      crudy={GalleryCrudy}
      pageSize={10}
      searchParams={sp}
      searchPropName="like_name"
      inKeyword="in_id"
    />
  );
}
