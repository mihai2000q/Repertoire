export default function fixDraggableInMenu(
  draggableProps,
  snapshot: { isDragging: boolean },
  offset = 160
) {
  if (snapshot.isDragging) {
    if ('left' in draggableProps.style) {
      draggableProps.style.left = 0
    }
    if ('top' in draggableProps.style) {
      draggableProps.style.top -= offset
    }
  }

  return draggableProps
}
