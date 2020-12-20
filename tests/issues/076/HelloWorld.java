/*
 * Copyright (C) 2014, 2020 Yoshiki Shibata. All rights reserved.
 */
import javafx.application.Application;
import javafx.scene.Scene;
import javafx.scene.control.Label;
import javafx.scene.control.TextArea;
import javafx.scene.layout.VBox;
import javafx.scene.text.Font;
import javafx.stage.Stage;

/**
 *
 * Write a program with a text field and a label. As with the Hello, JavaFX
 * program, the label should have the string Hello, FX in a 100 point font.
 * Initialize the text field with the same string. Update the label as the user
 * edits the text field.
 *
 * テキストフィールドとラベルを持つプログラムを作成しなさい。「Hello, JavaFX」プログ ラムと同じように、そのラベルは、文字列Hello, FX
 * を100 ポイントのフォントで表 示するようにしなさい。テキストフィールドを同じ文字列で初期化しなさい。ユーザーが
 * テキストフィールドを編集したらラベルも更新するようにしなさい。
 */
public class HelloWorld extends Application {

    private static final String INITIAL_MESSAGE = "Hello, JavaFX!";

    @Override
    public void start(final Stage primaryStage) {

        TextArea input = new TextArea(INITIAL_MESSAGE);
        Label message = new Label(INITIAL_MESSAGE);
        message.setFont(new Font(100));
        message.textProperty().bindBidirectional(input.textProperty());
        VBox root = new VBox();
        root.getChildren().addAll(message, input);

        Scene scene = new Scene(root);

        primaryStage.setTitle("Hello World!");
        primaryStage.setScene(scene);
        primaryStage.show();
    }
}
