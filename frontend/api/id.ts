import { useUserStore } from "@/stores/user"
import type { InitMessage } from "./interfaces"

const idLSKey = "id"

const handlePacket = (data: InitMessage) => {
  localStorage.setItem(idLSKey, data.id ?? "wrong uuid")

  const user = useUserStore()
  user.data = data
}

const getLocalId = () => {
  return localStorage.getItem(idLSKey)
}

export default {handlePacket, getLocalId}