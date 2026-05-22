interface SelectItem {
  label: string
  value: string
}

type GroupedItemList = Array<{
  label: string
  items: Array<SelectItem>
}>
