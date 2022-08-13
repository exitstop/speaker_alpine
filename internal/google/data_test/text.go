package data_test

var (
	Text = []string{
		`

1508
[["wrb.fr","MkEWBc","[[null,null,\"en\",[[[0,[[[null,35]],[true]]],[1,[[[null,36],[36,177]],[false,true]]]],177],[[\"is changed without server requests.\",null,null,35],[\" \",null,35,36],[\"For example, when you click a button and something pops on the screen, or some content on the page changes, etc... i think you take the idea.\",null,36,177]]],[[[null,\"izmenyayetsya bez zaprosov servera. Naprimer, kogda vy nazhimayete knopku i chto -to poyavlyayetsya na ekrane, ili kakoy -to kontent na stranitse menyayetsya i t. D. YA dumayu, chto vy prinimayete etu ideyu.\",null,null,null,[[\"изменяется без запросов сервера.\",null,null,null,[[\"изменяется без запросов сервера.\",[5]],[\"изменяется без запросов на сервера.\",[11]]]],[\"Например, когда вы нажимаете кнопку и что -то появляется на экране, или какой -то контент на странице меняется и т. Д. Я думаю, что вы принимаете эту идею.\",null,true,null,[[\"Например, когда вы нажимаете кнопку и что -то появляется на экране, или какой -то контент на странице меняется и т. Д. Я думаю, что вы принимаете эту идею.\",[5]],[\"Например, когда вы нажимаете кнопку и что -то появляется на экране, или какой -то контент на странице изменяется и т. Д. Я думаю, что вы принимаете эту идею.\",[11]]]]]]],\"ru\",1,\"en\",[\"is changed without server requests. For example, when you click a button and something pops on the screen, or some content on the page changes, etc... i think you take the idea.\",\"auto\",\"ru\",true]],\"en\"]",null,null,null,"generic"]]
58
[["di",482],["af.httprm",481,"-3052799092954120545",11]]
26
[["e",4,null,null,2054]]

`,
		`

1069
[["wrb.fr","MkEWBc","[[null,null,\"en\",[[[0,[[[null,113]],[true]]]],113],[[\"A static page is a page delivered to the user exactly as stored and with no chance on being changed, end of story\",null,null,113]]],[[[null,\"Staticheskaya stranitsa - eto stranitsa, dostavlennaya pol'zovatelyu tochno tak zhe, kak khranitsya, i bez shansov na izmeneniye, konets istorii\",null,null,null,[[\"Статическая страница - это страница, доставленная пользователю точно так же, как хранится, и без шансов на изменение, конец истории\",null,null,null,[[\"Статическая страница - это страница, доставленная пользователю точно так же, как хранится, и без шансов на изменение, конец истории\",[5]],[\"Статическая страница - это страница, доставленная пользователю точно так же, как хранится, и без шансов изменить, конец истории\",[11]]]]]]],\"ru\",1,\"en\",[\"A static page is a page delivered to the user exactly as stored and with no chance on being changed, end of story\",\"auto\",\"ru\",true]],\"en\"]",null,null,null,"generic"],["di",31],["af.httprm",30,"9058989769143286083",8]]
26
[["e",4,null,null,1428]]

`,
	}

	TextOut = []string{
		"изменяется без запросов сервера. Например, когда вы нажимаете кнопку и что -то появляется на экране, или какой -то контент на странице меняется и т. Д. Я думаю, что вы принимаете эту идею.",
	}
)
