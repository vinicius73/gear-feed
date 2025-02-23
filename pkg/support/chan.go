package support

import "sync"

// MergeChannels combina múltiplos canais de entrada em um único canal de saída.
// Os valores de todos os canais de entrada são enviados para o canal de saída.
// O canal de saída é fechado quando todos os canais de entrada forem fechados.
func MergeChannels[T any](list ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	// Adiciona um buffer ao canal de saída para melhorar o desempenho
	out := make(chan T, 10) // O tamanho do buffer pode ser ajustado conforme necessário

	output := func(c <-chan T) {
		defer wg.Done()
		// Verifica se o canal é nil para evitar bloqueios
		if c == nil {
			return
		}
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(list))
	for _, c := range list {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
