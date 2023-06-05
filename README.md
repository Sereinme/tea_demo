# TEA-DEMO

My first bubble tea demo using go.

![](.\demo.gif)

## NOTE

When using `lipgloss` to customize the style of std output, we use `Render` method to change the strings. But if we use `Render` to `"\n"`, the std will get extra spaces that will change the whole tui performance. So do not use `Render` to `"\n"`.
