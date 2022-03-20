package com.gusrylmubarok.springdi;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

class SpringDiApplicationTests {

	public static void main(String[] args) {
		ApplicationContext applicationContext = new AnnotationConfigApplicationContext(SpringDiApplication.class);
		SetterBasedInjection  fieldBasedInjection = applicationContext.getBean(SetterBasedInjection.class);
		fieldBasedInjection.processMsg("twitter message sending ");
	}
}
