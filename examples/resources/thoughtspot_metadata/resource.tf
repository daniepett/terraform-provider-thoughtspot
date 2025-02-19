resource "thoughtspot_tml" "this" {
  tml = <<EOT
guid: d084c256-e284-4fc4-b80c-111cb6064300
liveboard:
  name: Sales Performance
  visualizations:
  - id: Viz_1
    answer:
      name: Total sales by store
      tables:
      - id: (Sample) Retail - Apparel
        name: (Sample) Retail - Apparel
      search_query: "[sales] [store] [date].2021 top 10"
      answer_columns:
      - name: Total sales
      - name: store
      table:
        table_columns:
        - column_id: Total sales
          headline_aggregation: SUM
        - column_id: store
          headline_aggregation: COUNT_DISTINCT
        ordered_column_ids:
        - store
        - Total sales
      chart:
        type: COLUMN
        chart_columns:
        - column_id: Total sales
        - column_id: store
        axis_configs:
        - x:
          - store
          "y":
          - Total sales
      display_mode: CHART_MODE
    viz_guid: 5adc64e6-e631-4bf9-bd52-10ce17195300
  layout:
    tiles:
    - visualization_id: Viz_1
      size: MEDIUM_SMALL

EOT
}