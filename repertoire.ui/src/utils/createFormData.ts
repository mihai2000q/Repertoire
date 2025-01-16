export default function createFormData(obj: object) {
  const formData = new FormData()
  Object.entries(obj).forEach(([key, value]) => formData.append(key, value))
  return formData
}
