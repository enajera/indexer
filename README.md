# Indexer :mailbox_with_mail:
Programa en Go que indexa correos electronicos a Zincsearch

![Process](internal/process/Process.png)

## Profiling Mejorado aplicando concurrencia 📊
### CPU - Heap 

#### Antes
Archivos procesados: 517425

Tiempo de ejecución: 2h37m23.0416788s
![Antes](internal/profiling/profile_heap.png)


#### Despues
Archivos procesados: 517425

Tiempo Total: 37m59.3782087s 
![Despues](internal/profiling/profile_mejorado.png)
