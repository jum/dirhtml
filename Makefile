TARG=dirhtml
GOFILES=dirhtml.go

$(TARG): $(GOFILES)
	go build -o $(TARG) $(GOFILES)

clean:
	rm -f $(TARG) mon.out
